package ast

import (
	"fmt"

	"github.com/textwire/textwire/token"
)

type StringLiteral struct {
	Token token.Token
	Value string
}

func (i *StringLiteral) expressionNode() {
}

func (i *StringLiteral) TokenLiteral() string {
	return i.Token.Literal
}

func (i *StringLiteral) String() string {
	return fmt.Sprintf(`"%s"`, i.Value)
}
