package ast

import (
	"strings"

	"github.com/textwire/textwire/v3/pkg/token"
)

type Embedded struct {
	BaseNode
	Elements []EmbeddedElement
}

func NewEmbedded(tok token.Token) *Embedded {
	return &Embedded{
		BaseNode: NewBaseNode(tok),
	}
}

func (_ *Embedded) chunkNode() {}

func (_ *Embedded) Kind() ChunkKind {
	return ChunkKindEmbedded
}

func (e *Embedded) String() string {
	var out strings.Builder
	out.Grow(4)

	out.WriteString("{{ ")

	for i, stmt := range e.Elements {
		out.WriteString(stmt.String())

		if i < len(e.Elements)-1 {
			out.WriteString("; ")
		}
	}

	out.WriteString(" }}")

	return out.String()
}
