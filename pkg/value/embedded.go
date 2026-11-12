package value

import (
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

	toString := escapedString
	if e.IsRaw {
		toString = Literal.String
	}

	for i := range e.Segments {
		out.WriteString(toString(e.Segments[i]))
	}

	return out.String()
}

func (e *Embedded) Is(t ValueType) bool {
	return t == e.Type()
}
