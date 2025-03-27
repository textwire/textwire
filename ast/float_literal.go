package ast

import (
	"fmt"

	"github.com/textwire/textwire/v2/token"
)

type FloatLiteral struct {
	Token token.Token
	Value float64
	Pos   token.Position
}

func (fl *FloatLiteral) expressionNode() {
}

func (fl *FloatLiteral) Tok() *token.Token {
	return &fl.Token
}

func (fl *FloatLiteral) String() string {
	return fmt.Sprintf("%g", fl.Value)
}

func (fl *FloatLiteral) Line() uint {
	return fl.Token.ErrorLine()
}

func (fl *FloatLiteral) Position() token.Position {
	return fl.Pos
}
