package htmlconv

import (
	"testing"
)

func TestConvertHTML(t *testing.T) {
	input := []byte(`<html><body>
<h1>Title</h1>
<p>Hello <b>bold</b> world</p>
<p>Second paragraph with <i>italic</i> text</p>
</body></html>`)

	lines, err := ConvertHTML(input)
	if err != nil {
		t.Fatalf("ConvertHTML failed: %v", err)
	}

	if len(lines) == 0 {
		t.Fatal("expected lines, got none")
	}

	foundBold := false
	foundItalic := false
	for _, line := range lines {
		for _, span := range line.Spans {
			if span.Style.Bold && span.Text == "Title" {
				foundBold = true
			}
			if span.Style.Italic && span.Text == "italic" {
				foundItalic = true
			}
		}
	}

	if !foundBold {
		t.Error("did not find bold 'Title' span")
	}
	if !foundItalic {
		t.Error("did not find italic span")
	}
}

func TestScriptAndStyleExcluded(t *testing.T) {
	input := []byte(`<html><body>
<script>var x = 1;</script>
<style>.foo { color: red; }</style>
<p>Visible text</p>
</body></html>`)

	lines, err := ConvertHTML(input)
	if err != nil {
		t.Fatalf("ConvertHTML failed: %v", err)
	}

	for _, line := range lines {
		for _, span := range line.Spans {
			if span.Text == "var x = 1;" || span.Text == ".foo { color: red; }" {
				t.Error("script/style content should be excluded")
			}
		}
	}
}
