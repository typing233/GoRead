package ui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"

	"github.com/goread/goread/internal/config"
	"github.com/goread/goread/internal/htmlconv"
	"github.com/goread/goread/internal/paginator"
)

type Renderer struct {
	screen    tcell.Screen
	baseStyle tcell.Style
	cfg       *config.Config
}

func NewRenderer(screen tcell.Screen, cfg *config.Config) *Renderer {
	r := &Renderer{
		screen: screen,
		cfg:    cfg,
	}
	r.baseStyle = r.buildBaseStyle()
	return r
}

func (r *Renderer) buildBaseStyle() tcell.Style {
	style := tcell.StyleDefault

	fg := parseColor(r.cfg.ForegroundColor)
	bg := parseColor(r.cfg.BackgroundColor)
	style = style.Foreground(fg).Background(bg)

	if r.cfg.Bold {
		style = style.Bold(true)
	}
	if r.cfg.Italic {
		style = style.Italic(true)
	}
	if r.cfg.Underline {
		style = style.Underline(true)
	}

	return style
}

func (r *Renderer) RenderPage(lines []paginator.WrappedLine, offset, totalLines int, chapterTitle string, chapterIdx, totalChapters int) {
	r.screen.Clear()
	width, height := r.screen.Size()
	margin := r.cfg.SideMargin
	contentWidth := width - 2*margin
	if contentWidth < 10 {
		margin = 0
		contentWidth = width
	}

	pageHeight := height - 1

	for i, line := range lines {
		if i >= pageHeight {
			break
		}
		col := margin
		for _, span := range line.Spans {
			style := r.applySpanStyle(span.Style)
			for _, ch := range span.Text {
				rw := runewidth.RuneWidth(ch)
				if col+rw > width-margin {
					break
				}
				r.screen.SetContent(col, i, ch, nil, style)
				col += rw
			}
		}
	}

	r.renderStatusBar(width, height, offset, totalLines, chapterTitle, chapterIdx, totalChapters)
	r.screen.Show()
}

func (r *Renderer) renderStatusBar(width, height, offset, totalLines int, chapterTitle string, chapterIdx, totalChapters int) {
	y := height - 1
	statusStyle := tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorWhite)

	for x := 0; x < width; x++ {
		r.screen.SetContent(x, y, ' ', nil, statusStyle)
	}

	var pct int
	if totalLines > 0 {
		pct = (offset * 100) / totalLines
	}
	left := fmt.Sprintf(" [%d/%d] %s", chapterIdx+1, totalChapters, chapterTitle)
	right := fmt.Sprintf("%d%% ", pct)

	for i, ch := range left {
		if i >= width-len(right)-1 {
			break
		}
		r.screen.SetContent(i, y, ch, nil, statusStyle)
	}

	startX := width - len(right)
	for i, ch := range right {
		r.screen.SetContent(startX+i, y, ch, nil, statusStyle)
	}
}

func (r *Renderer) applySpanStyle(s htmlconv.TextStyle) tcell.Style {
	style := r.baseStyle
	if s.Bold || s.Heading > 0 {
		style = style.Bold(true)
	}
	if s.Italic {
		style = style.Italic(true)
	}
	if s.Underline {
		style = style.Underline(true)
	}
	return style
}

func parseColor(color string) tcell.Color {
	color = strings.TrimSpace(strings.ToLower(color))
	if color == "" || color == "default" {
		return tcell.ColorDefault
	}

	if strings.HasPrefix(color, "#") && len(color) == 7 {
		r, err1 := strconv.ParseInt(color[1:3], 16, 32)
		g, err2 := strconv.ParseInt(color[3:5], 16, 32)
		b, err3 := strconv.ParseInt(color[5:7], 16, 32)
		if err1 == nil && err2 == nil && err3 == nil {
			return tcell.NewRGBColor(int32(r), int32(g), int32(b))
		}
	}

	namedColors := map[string]tcell.Color{
		"black":   tcell.ColorBlack,
		"red":     tcell.ColorRed,
		"green":   tcell.ColorGreen,
		"yellow":  tcell.ColorYellow,
		"blue":    tcell.ColorBlue,
		"magenta": tcell.ColorDarkMagenta,
		"cyan":    tcell.ColorDarkCyan,
		"white":   tcell.ColorWhite,
	}

	if c, ok := namedColors[color]; ok {
		return c
	}
	return tcell.ColorDefault
}
