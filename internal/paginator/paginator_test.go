package paginator

import (
	"testing"

	"github.com/goread/goread/internal/htmlconv"
)

func TestWrapLine(t *testing.T) {
	line := htmlconv.StyledLine{
		Spans: []htmlconv.StyledSpan{
			{Text: "Hello world this is a long line that should wrap", Style: htmlconv.TextStyle{}},
		},
	}

	wrapped := wrapLine(line, 20)
	if len(wrapped) < 2 {
		t.Errorf("expected multiple wrapped lines, got %d", len(wrapped))
	}
}

func TestPaginator(t *testing.T) {
	lines := []htmlconv.StyledLine{
		{Spans: []htmlconv.StyledSpan{{Text: "Line 1"}}},
		{Spans: []htmlconv.StyledSpan{{Text: "Line 2"}}},
		{Spans: []htmlconv.StyledSpan{{Text: "Line 3"}}},
		{Spans: []htmlconv.StyledSpan{{Text: "Line 4"}}},
		{Spans: []htmlconv.StyledSpan{{Text: "Line 5"}}},
	}

	p := New(lines, 80, 3)
	if p.TotalLines() != 5 {
		t.Errorf("expected 5 total lines, got %d", p.TotalLines())
	}

	page := p.GetLines(0, 3)
	if len(page) != 3 {
		t.Errorf("expected 3 lines on page, got %d", len(page))
	}

	page = p.GetLines(3, 3)
	if len(page) != 2 {
		t.Errorf("expected 2 lines on last page, got %d", len(page))
	}
}
