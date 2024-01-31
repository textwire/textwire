package ast

import (
	"fmt"

	"github.com/textwire/textwire/token"
)

type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (fl *FloatLiteral) expressionNode() {
}

func (fl *FloatLiteral) TokenLiteral() string {
	return fl.Token.Literal
}

func (fl *FloatLiteral) String() string {
	return fmt.Sprintf("%g", fl.Value)
}

func (fl *FloatLiteral) LineNum() uint {
	return fl.Token.Line
}
