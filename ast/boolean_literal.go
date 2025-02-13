package ast

import (
	"github.com/textwire/textwire/v2/token"
)

type BooleanLiteral struct {
	Token token.Token // The 'true' or 'false' token
	Value bool
	Pos   Position
}

func (bl *BooleanLiteral) expressionNode() {
}

func (bl *BooleanLiteral) TokenLiteral() string {
	return bl.Token.Literal
}

func (bl *BooleanLiteral) String() string {
	return bl.Token.Literal
}

func (bl *BooleanLiteral) Line() uint {
	return bl.Token.DebugLine
}

func (bl *BooleanLiteral) Position() Position {
	return bl.Pos
}
