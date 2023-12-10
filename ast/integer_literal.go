package ast

import "github.com/textwire/textwire/token"

type IntegerLiteral struct {
	Token token.Token
	Value int
}

func (i *IntegerLiteral) expressionNode() {
}

func (i *IntegerLiteral) TokenLiteral() string {
	return i.Token.Literal
}

func (i *IntegerLiteral) String() string {
	return i.Token.Literal
}
