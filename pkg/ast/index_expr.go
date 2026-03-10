package ast

import (
	"fmt"

	"github.com/textwire/textwire/v3/pkg/token"
)

type IndexExpr struct {
	BaseNode
	Left  Expression
	Index Expression
}

func NewIndexExpr(tok token.Token, left Expression) *IndexExpr {
	return &IndexExpr{
		BaseNode: NewBaseNode(tok),
		Left:     left,
	}
}

func (*IndexExpr) expressionNode() {}

func (ie *IndexExpr) String() string {
	return fmt.Sprintf("(%s[%s])", ie.Left, ie.Index)
}
