package value

import (
	"strings"
)

type Embedded struct {
	Segments []Literal
}

func NewEmbedded(cap int) *Embedded {
	return &Embedded{Segments: make([]Literal, 0, cap)}
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

	return out.String()
}

func (b *Embedded) Is(t ValueType) bool {
	return t == b.Type()
}
