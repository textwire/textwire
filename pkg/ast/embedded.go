package ast

import (
	"html"
	"strings"

	"github.com/textwire/textwire/v5/pkg/token"
)

type Embedded struct {
	BaseNode
	Segments []Segment
	IsRaw    bool
}

func NewEmbedded(tok token.Token) *Embedded {
	return &Embedded{
		BaseNode: NewBaseNode(tok),
	}
}

func (*Embedded) chunkNode() {}

func (e *Embedded) String() string {
	var out strings.Builder
	out.Grow(4)

	out.WriteString("{{ ")

	for i, stmt := range e.Segments {
		out.WriteString(stmt.String())

		if i < len(e.Segments)-1 {
			out.WriteString("; ")
		}
	}

	out.WriteString(" }}")

	if e.IsRaw {
		return out.String()
	}
	return html.EscapeString(out.String())
}
