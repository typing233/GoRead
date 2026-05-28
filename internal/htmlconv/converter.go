package htmlconv

import (
	"bytes"
	"strings"

	"golang.org/x/net/html"
)

func ConvertHTML(xhtml []byte) ([]StyledLine, error) {
	doc, err := html.Parse(bytes.NewReader(xhtml))
	if err != nil {
		return nil, err
	}

	c := &converter{}
	c.walk(doc, TextStyle{})
	c.flush()
	return c.lines, nil
}

type converter struct {
	lines   []StyledLine
	current []StyledSpan
	buf     strings.Builder
	style   TextStyle
}

func (c *converter) flush() {
	if c.buf.Len() > 0 {
		c.current = append(c.current, StyledSpan{
			Text:  c.buf.String(),
			Style: c.style,
		})
		c.buf.Reset()
	}
	if len(c.current) > 0 {
		c.lines = append(c.lines, StyledLine{Spans: c.current})
		c.current = nil
	}
}

func (c *converter) flushSpan() {
	if c.buf.Len() > 0 {
		c.current = append(c.current, StyledSpan{
			Text:  c.buf.String(),
			Style: c.style,
		})
		c.buf.Reset()
	}
}

func (c *converter) emitLineBreak() {
	if c.buf.Len() > 0 {
		c.current = append(c.current, StyledSpan{
			Text:  c.buf.String(),
			Style: c.style,
		})
		c.buf.Reset()
	}
	c.lines = append(c.lines, StyledLine{Spans: c.current})
	c.current = nil
}

func (c *converter) walk(n *html.Node, parentStyle TextStyle) {
	switch n.Type {
	case html.TextNode:
		text := collapseWhitespace(n.Data)
		if text != "" {
			c.buf.WriteString(text)
		}
		return
	case html.ElementNode:
		tag := strings.ToLower(n.Data)
		newStyle := parentStyle

		switch tag {
		case "b", "strong":
			newStyle.Bold = true
		case "i", "em", "cite":
			newStyle.Italic = true
		case "u":
			newStyle.Underline = true
		case "h1":
			newStyle.Bold = true
			newStyle.Heading = 1
		case "h2":
			newStyle.Bold = true
			newStyle.Heading = 2
		case "h3":
			newStyle.Bold = true
			newStyle.Heading = 3
		case "h4", "h5", "h6":
			newStyle.Bold = true
			newStyle.Heading = 4
		case "script", "style", "head":
			return
		case "br":
			c.emitLineBreak()
			return
		}

		isBlock := isBlockElement(tag)
		if isBlock {
			c.emitLineBreak()
		}

		if newStyle != c.style {
			c.flushSpan()
			oldStyle := c.style
			c.style = newStyle
			for child := n.FirstChild; child != nil; child = child.NextSibling {
				c.walk(child, newStyle)
			}
			c.flushSpan()
			c.style = oldStyle
		} else {
			for child := n.FirstChild; child != nil; child = child.NextSibling {
				c.walk(child, newStyle)
			}
		}

		if isBlock {
			c.emitLineBreak()
		}
		return
	}

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		c.walk(child, parentStyle)
	}
}

func isBlockElement(tag string) bool {
	switch tag {
	case "p", "div", "section", "article", "header", "footer",
		"h1", "h2", "h3", "h4", "h5", "h6",
		"ul", "ol", "li", "blockquote", "pre", "hr", "table",
		"tr", "nav", "main", "aside", "figure", "figcaption":
		return true
	}
	return false
}

func collapseWhitespace(s string) string {
	var buf strings.Builder
	inSpace := false
	for _, r := range s {
		if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
			if !inSpace {
				buf.WriteByte(' ')
				inSpace = true
			}
		} else {
			buf.WriteRune(r)
			inSpace = false
		}
	}
	return buf.String()
}
