package ast

import "github.com/textwire/textwire/v2/token"

type NilLiteral struct {
	Token token.Token // The 'nil' token
	Pos   Position
}

func (nl *NilLiteral) expressionNode() {
}

func (nl *NilLiteral) TokenLiteral() string {
	return nl.Token.Literal
}

func (nl *NilLiteral) String() string {
	return nl.Token.Literal
}

func (nl *NilLiteral) Line() uint {
	return nl.Token.DebugLine
}

func (nl *NilLiteral) Position() Position {
	return nl.Pos
}
