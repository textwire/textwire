package ast

import "github.com/textwire/textwire/v2/token"

type IntegerLiteral struct {
	Token token.Token
	Value int64
	Pos   Position
}

func (il *IntegerLiteral) expressionNode() {
}

func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}

func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}

func (il *IntegerLiteral) Line() uint {
	return il.Token.StartLine
}

func (il *IntegerLiteral) Position() Position {
	return il.Pos
}
