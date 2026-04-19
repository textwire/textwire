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

func (b *Embedded) String() string {
	var out strings.Builder
	out.Grow(len(b.Segments))

	for i := range b.Segments {
		out.WriteString(b.Segments[i].String())
	}

	if b.IsRaw {
		return out.String()
	}
	return html.EscapeString(out.String())
}

func (b *Embedded) Is(t ValueType) bool {
	return t == b.Type()
}
