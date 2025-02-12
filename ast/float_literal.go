package ast

import (
	"fmt"

	"github.com/textwire/textwire/v2/token"
)

type FloatLiteral struct {
	Token token.Token
	Value float64
	Pos   Position
}

func (fl *FloatLiteral) expressionNode() {
}

func (fl *FloatLiteral) TokenLiteral() string {
	return fl.Token.Literal
}

func (fl *FloatLiteral) String() string {
	return fmt.Sprintf("%g", fl.Value)
}

func (fl *FloatLiteral) Line() uint {
	return fl.Token.StartLine
}

func (fl *FloatLiteral) Position() Position {
	return fl.Pos
}
