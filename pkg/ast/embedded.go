package ast

import (
	"strings"

	"github.com/textwire/textwire/v5/pkg/token"
)

type Embedded struct {
	BaseNode
	Segments []Segment
	IsRaw    bool
}

func NewEmbedded(tok token.Token, isRaw bool) *Embedded {
	return &Embedded{
		BaseNode: NewBaseNode(tok),
		IsRaw:    isRaw,
	}
}

func (*Embedded) chunkNode() {}

func (e *Embedded) String() string {
	var out strings.Builder
	out.Grow(4)

	if e.IsRaw {
		out.WriteString("{!! ")
	} else {
		out.WriteString("{{ ")
	}

	for i, stmt := range e.Segments {
		out.WriteString(stmt.String())

		if i < len(e.Segments)-1 {
			out.WriteString("; ")
		}
	}

	if e.IsRaw {
		out.WriteString(" !!}")
		return out.String()
	}

	out.WriteString(" }}")
	return out.String()
}
