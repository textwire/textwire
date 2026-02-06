package ast

import (
	"strings"

	"github.com/textwire/textwire/v3/token"
)

type ObjectLiteral struct {
	BaseNode
	Pairs map[string]Expression // Key-value pairs; { key: value }
}

func NewObjectLiteral(tok token.Token) *ObjectLiteral {
	return &ObjectLiteral{
		BaseNode: NewBaseNode(tok),
	}
}

func (ol *ObjectLiteral) expressionNode() {}

func (ol *ObjectLiteral) String() string {
	var out strings.Builder
	estimateSize := 2 + len(ol.Pairs)
	estimateSize += 2 + len(ol.Pairs) // quotes
	estimateSize += 2 + len(ol.Pairs) // colon with space
	out.Grow(estimateSize)

	out.WriteString("{")

	for key, value := range ol.Pairs {
		out.WriteByte('"')
		out.WriteString(key)
		out.WriteString(`": `)
		out.WriteString(value.String())
	}

	out.WriteString("}")

	return out.String()

}
