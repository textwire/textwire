package ast

import (
	"fmt"

	"github.com/textwire/textwire/token"
)

type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (f *FloatLiteral) expressionNode() {
}

func (f *FloatLiteral) TokenLiteral() string {
	return f.Token.Literal
}

func (f *FloatLiteral) String() string {
	return fmt.Sprintf("%g", f.Value)
}
