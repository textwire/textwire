package ast

import (
	"fmt"

	"github.com/textwire/textwire/v3/pkg/token"
)

type DotExpr struct {
	BaseNode
	Left Expression // -->x.y
	Key  Expression // x.y<--
}

func NewDotExpr(tok token.Token, left Expression) *DotExpr {
	return &DotExpr{
		BaseNode: NewBaseNode(tok),
		Left:     left,
	}
}

func (*DotExpr) expressionNode() {}
func (*DotExpr) segmentNode()    {}

func (de *DotExpr) String() string {
	return fmt.Sprintf("(%s.%s)", de.Left, de.Key)
}
