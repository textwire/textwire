package ast

import (
	"github.com/textwire/textwire/v2/token"
)

type BooleanLiteral struct {
	Token token.Token // The 'true' or 'false' token
	Value bool
	Pos   token.Position
}

func (bl *BooleanLiteral) expressionNode() {}

func (bl *BooleanLiteral) Tok() *token.Token {
	return &bl.Token
}

func (bl *BooleanLiteral) String() string {
	return bl.Token.Literal
}

func (bl *BooleanLiteral) Line() uint {
	return bl.Token.ErrorLine()
}

func (bl *BooleanLiteral) Position() token.Position {
	return bl.Pos
}
