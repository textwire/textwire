package ast

import (
	"fmt"

	"github.com/textwire/textwire/v3/pkg/token"
)

type StringLiteral struct {
	BaseNode
	Value string
}

func NewStringLiteral(tok token.Token, val string) *StringLiteral {
	return &StringLiteral{
		BaseNode: NewBaseNode(tok),
		Value:    val,
	}
}

func (sl *StringLiteral) expressionNode() {}

func (sl *StringLiteral) String() string {
	return fmt.Sprintf(`"%s"`, sl.Token.Literal)
}
