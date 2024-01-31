package ast

import "github.com/textwire/textwire/token"

type NilLiteral struct {
	Token token.Token // The 'nil' token
}

func (nl *NilLiteral) expressionNode() {
}

func (nl *NilLiteral) TokenLiteral() string {
	return nl.Token.Literal
}

func (nl *NilLiteral) String() string {
	return nl.Token.Literal
}

func (nl *NilLiteral) LineNum() uint {
	return nl.Token.Line
}
