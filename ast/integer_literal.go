package ast

import "github.com/textwire/textwire/v2/token"

type IntegerLiteral struct {
	Token token.Token
	Value int64
	Pos   token.Position
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
	return il.Token.ErrorLine()
}

func (il *IntegerLiteral) Position() token.Position {
	return il.Pos
}
