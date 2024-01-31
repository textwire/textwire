package ast

import (
	"github.com/textwire/textwire/token"
)

type StringLiteral struct {
	Token token.Token // The content of the string
	Value string
}

func (sl *StringLiteral) expressionNode() {
}

func (sl *StringLiteral) TokenLiteral() string {
	return sl.Token.Literal
}

func (sl *StringLiteral) String() string {
	return `"` + sl.Token.Literal + `"`
}

func (sl *StringLiteral) LineNum() uint {
	return sl.Token.Line
}
