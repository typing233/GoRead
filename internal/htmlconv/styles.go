package htmlconv

type TextStyle struct {
	Bold      bool
	Italic    bool
	Underline bool
	Heading   int
}

type StyledSpan struct {
	Text  string
	Style TextStyle
}

type StyledLine struct {
	Spans []StyledSpan
}
