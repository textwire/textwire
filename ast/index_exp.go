package ast

import (
	"fmt"

	"github.com/textwire/textwire/v3/token"
)

type IndexExp struct {
	BaseNode
	Left  Expression
	Index Expression
}

func NewIndexExp(tok token.Token, left Expression) *IndexExp {
	return &IndexExp{
		BaseNode: NewBaseNode(tok),
		Left:     left,
	}
}

func (ie *IndexExp) expressionNode() {}

func (ie *IndexExp) String() string {
	return fmt.Sprintf("(%s[%s])", ie.Left, ie.Index)
}
