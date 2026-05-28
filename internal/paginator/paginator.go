package paginator

import (
	"github.com/mattn/go-runewidth"

	"github.com/goread/goread/internal/htmlconv"
)

type WrappedLine struct {
	Spans []htmlconv.StyledSpan
}

type Paginator struct {
	lines        []htmlconv.StyledLine
	wrappedLines []WrappedLine
	width        int
	height       int
}

func New(lines []htmlconv.StyledLine, width, height int) *Paginator {
	p := &Paginator{
		lines:  lines,
		width:  width,
		height: height,
	}
	p.reflow()
	return p
}

func (p *Paginator) Reflow(width, height int) {
	p.width = width
	p.height = height
	p.reflow()
}

func (p *Paginator) TotalLines() int {
	return len(p.wrappedLines)
}

func (p *Paginator) GetLines(offset, count int) []WrappedLine {
	if offset < 0 {
		offset = 0
	}
	if offset >= len(p.wrappedLines) {
		return nil
	}
	end := offset + count
	if end > len(p.wrappedLines) {
		end = len(p.wrappedLines)
	}
	return p.wrappedLines[offset:end]
}

func (p *Paginator) reflow() {
	p.wrappedLines = nil
	for _, line := range p.lines {
		wrapped := wrapLine(line, p.width)
		p.wrappedLines = append(p.wrappedLines, wrapped...)
	}
	for len(p.wrappedLines) > 0 && isEmptyWrappedLine(p.wrappedLines[0]) {
		p.wrappedLines = p.wrappedLines[1:]
	}
}

func isEmptyWrappedLine(wl WrappedLine) bool {
	for _, span := range wl.Spans {
		for _, r := range span.Text {
			if r != ' ' && r != '\t' {
				return false
			}
		}
	}
	return true
}

func wrapLine(line htmlconv.StyledLine, width int) []WrappedLine {
	if width <= 0 {
		return []WrappedLine{{Spans: line.Spans}}
	}
	if len(line.Spans) == 0 {
		return []WrappedLine{{}}
	}

	var result []WrappedLine
	var currentSpans []htmlconv.StyledSpan
	col := 0

	for _, span := range line.Spans {
		runes := []rune(span.Text)
		pos := 0

		for pos < len(runes) {
			r := runes[pos]
			rw := runewidth.RuneWidth(r)

			if r == ' ' && col+rw > width {
				result = append(result, WrappedLine{Spans: currentSpans})
				currentSpans = nil
				col = 0
				pos++
				continue
			}

			if col+rw > width {
				// Current line is full — try to break at last space
				broke := breakCurrentLine(&result, &currentSpans, &col, width)
				if !broke {
					// No space found: force break at current position
					if len(currentSpans) > 0 || col > 0 {
						result = append(result, WrappedLine{Spans: currentSpans})
						currentSpans = nil
						col = 0
					}
				}
				continue
			}

			// Append this rune to current span accumulation
			currentSpans = appendRune(currentSpans, r, span.Style)
			col += rw
			pos++
		}
	}

	if len(currentSpans) > 0 || len(result) == 0 {
		result = append(result, WrappedLine{Spans: currentSpans})
	}
	return result
}

// breakCurrentLine tries to find the last space in currentSpans and splits there.
// Returns true if a break was made.
func breakCurrentLine(result *[]WrappedLine, currentSpans *[]htmlconv.StyledSpan, col *int, width int) bool {
	// Find last space position across all spans
	type breakPos struct {
		spanIdx int
		runeIdx int
	}
	var lastBreak *breakPos

	for si, sp := range *currentSpans {
		for ri, r := range []rune(sp.Text) {
			if r == ' ' {
				lastBreak = &breakPos{spanIdx: si, runeIdx: ri}
			}
		}
	}

	if lastBreak == nil {
		return false
	}

	// Split at the break point: everything before space goes to current line
	var beforeSpans []htmlconv.StyledSpan
	var afterSpans []htmlconv.StyledSpan

	spans := *currentSpans
	for si, sp := range spans {
		if si < lastBreak.spanIdx {
			beforeSpans = append(beforeSpans, sp)
		} else if si == lastBreak.spanIdx {
			runes := []rune(sp.Text)
			if lastBreak.runeIdx > 0 {
				beforeSpans = append(beforeSpans, htmlconv.StyledSpan{
					Text:  string(runes[:lastBreak.runeIdx]),
					Style: sp.Style,
				})
			}
			// Skip the space itself, take rest as after
			rest := lastBreak.runeIdx + 1
			if rest < len(runes) {
				afterSpans = append(afterSpans, htmlconv.StyledSpan{
					Text:  string(runes[rest:]),
					Style: sp.Style,
				})
			}
		} else {
			afterSpans = append(afterSpans, sp)
		}
	}

	*result = append(*result, WrappedLine{Spans: beforeSpans})
	*currentSpans = afterSpans
	*col = displayWidth(afterSpans)
	return true
}

func appendRune(spans []htmlconv.StyledSpan, r rune, style htmlconv.TextStyle) []htmlconv.StyledSpan {
	if len(spans) > 0 && spans[len(spans)-1].Style == style {
		spans[len(spans)-1].Text += string(r)
	} else {
		spans = append(spans, htmlconv.StyledSpan{
			Text:  string(r),
			Style: style,
		})
	}
	return spans
}

func displayWidth(spans []htmlconv.StyledSpan) int {
	w := 0
	for _, sp := range spans {
		w += runewidth.StringWidth(sp.Text)
	}
	return w
}
