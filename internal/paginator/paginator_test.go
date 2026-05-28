package paginator

import (
	"testing"

	"github.com/mattn/go-runewidth"

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

	// Verify no wrapped line exceeds width
	for i, wl := range wrapped {
		w := 0
		for _, sp := range wl.Spans {
			w += runewidth.StringWidth(sp.Text)
		}
		if w > 20 {
			t.Errorf("wrapped line %d has display width %d, exceeds limit 20", i, w)
		}
	}
}

func TestWrapCJK(t *testing.T) {
	// Each CJK char is 2 columns wide. "你好世界" = 8 columns
	line := htmlconv.StyledLine{
		Spans: []htmlconv.StyledSpan{
			{Text: "你好世界测试文本换行", Style: htmlconv.TextStyle{}},
		},
	}

	// Width 10 = can fit 5 CJK chars per line. "你好世界测试文本换行" = 10 chars, 20 cols → 2 lines
	wrapped := wrapLine(line, 10)
	if len(wrapped) < 2 {
		t.Errorf("expected CJK text to wrap, got %d lines", len(wrapped))
	}

	for i, wl := range wrapped {
		w := 0
		for _, sp := range wl.Spans {
			w += runewidth.StringWidth(sp.Text)
		}
		if w > 10 {
			t.Errorf("CJK wrapped line %d has display width %d, exceeds limit 10", i, w)
		}
	}
}

func TestWrapMixed(t *testing.T) {
	// Mix of ASCII and CJK: "Hello你好World" = 5 + 4 + 5 = 14 cols
	line := htmlconv.StyledLine{
		Spans: []htmlconv.StyledSpan{
			{Text: "Hello你好World", Style: htmlconv.TextStyle{}},
		},
	}

	wrapped := wrapLine(line, 10)
	for i, wl := range wrapped {
		w := 0
		for _, sp := range wl.Spans {
			w += runewidth.StringWidth(sp.Text)
		}
		if w > 10 {
			t.Errorf("mixed wrapped line %d has display width %d, exceeds limit 10", i, w)
		}
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
