package ast

import "github.com/textwire/textwire/token"

type NilLiteral struct {
	Token token.Token
}

func (i *NilLiteral) expressionNode() {
}

func (n *NilLiteral) TokenLiteral() string {
	return n.Token.Literal
}

func (n *NilLiteral) String() string {
	return n.Token.Literal
}
