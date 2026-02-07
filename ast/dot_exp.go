package ast

import (
	"fmt"

	"github.com/textwire/textwire/v3/token"
)

type DotExp struct {
	BaseNode
	Left Expression // -->x.y
	Key  Expression // x.y<--
}

func NewDotExp(tok token.Token, left Expression) *DotExp {
	return &DotExp{
		BaseNode: NewBaseNode(tok),
		Left:     left,
	}
}

func (de *DotExp) expressionNode() {}

func (de *DotExp) String() string {
	return fmt.Sprintf("(%s.%s)", de.Left, de.Key)
}
