package value

import (
	"html"
	"strings"
)

type Embedded struct {
	Segments []Literal
	IsRaw    bool
}

func NewEmbedded(cap int, isRaw bool) *Embedded {
	return &Embedded{Segments: make([]Literal, 0, cap), IsRaw: isRaw}
}

func (*Embedded) Type() ValueType {
	return EMBEDDED_VAL
}

func (e *Embedded) String() string {
	var out strings.Builder
	out.Grow(len(e.Segments))

	for i := range e.Segments {
		out.WriteString(e.Segments[i].String())
	}

	if e.IsRaw {
		return out.String()
	}
	return html.EscapeString(out.String())
}

func (e *Embedded) Is(t ValueType) bool {
	return t == e.Type()
}
