package ast

import (
	"strings"

	"github.com/textwire/textwire/v3/pkg/token"
)

type ObjExpr struct {
	BaseNode
	Pairs map[string]Expression // Key-value pairs; { key: value }
}

func NewObjExpr(tok token.Token) *ObjExpr {
	return &ObjExpr{
		BaseNode: NewBaseNode(tok),
	}
}

func (*ObjExpr) expressionNode() {}

func (oe *ObjExpr) String() string {
	var out strings.Builder
	estimateSize := 2 + len(oe.Pairs)
	estimateSize += 2 + len(oe.Pairs) // quotes
	estimateSize += 2 + len(oe.Pairs) // colon with space
	out.Grow(estimateSize)

	out.WriteString("{")

	for key, value := range oe.Pairs {
		out.WriteByte('"')
		out.WriteString(key)
		out.WriteString(`": `)
		out.WriteString(value.String())
	}

	out.WriteString("}")

	return out.String()

}
