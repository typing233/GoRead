package paginator

import (
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
	// Trim leading empty lines
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
		text := span.Text
		for len(text) > 0 {
			remaining := width - col
			if remaining <= 0 {
				result = append(result, WrappedLine{Spans: currentSpans})
				currentSpans = nil
				col = 0
				remaining = width
			}

			if runeLen(text) <= remaining {
				currentSpans = append(currentSpans, htmlconv.StyledSpan{
					Text:  text,
					Style: span.Style,
				})
				col += runeLen(text)
				text = ""
			} else {
				breakAt := findBreakPoint(text, remaining)
				if breakAt <= 0 {
					if col == 0 {
						breakAt = remaining
					} else {
						result = append(result, WrappedLine{Spans: currentSpans})
						currentSpans = nil
						col = 0
						continue
					}
				}
				chunk := substringRunes(text, 0, breakAt)
				currentSpans = append(currentSpans, htmlconv.StyledSpan{
					Text:  chunk,
					Style: span.Style,
				})
				result = append(result, WrappedLine{Spans: currentSpans})
				currentSpans = nil
				col = 0
				text = substringRunes(text, breakAt, runeLen(text))
				text = trimLeadingSpace(text)
			}
		}
	}

	if len(currentSpans) > 0 || len(result) == 0 {
		result = append(result, WrappedLine{Spans: currentSpans})
	}

	return result
}

func findBreakPoint(text string, maxWidth int) int {
	lastSpace := -1
	col := 0
	for i, r := range text {
		if col >= maxWidth {
			break
		}
		if r == ' ' {
			lastSpace = col
			_ = i
		}
		col++
	}
	if lastSpace > 0 {
		return lastSpace
	}
	return -1
}

func runeLen(s string) int {
	n := 0
	for range s {
		n++
	}
	return n
}

func substringRunes(s string, start, end int) string {
	runes := []rune(s)
	if start > len(runes) {
		start = len(runes)
	}
	if end > len(runes) {
		end = len(runes)
	}
	return string(runes[start:end])
}

func trimLeadingSpace(s string) string {
	for i, r := range s {
		if r != ' ' {
			return s[i:]
		}
	}
	return ""
}
