package ast

import (
	"strings"

	"github.com/textwire/textwire/v3/pkg/token"
)

type ObjLit struct {
	BaseNode
	Pairs map[string]Expression // Key-value pairs; { key: value }
}

func NewObjLit(tok token.Token) *ObjLit {
	return &ObjLit{
		BaseNode: NewBaseNode(tok),
	}
}

func (ol *ObjLit) expressionNode() {}

func (ol *ObjLit) String() string {
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
