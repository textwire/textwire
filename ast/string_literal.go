package ast

import "github.com/textwire/textwire/v2/token"

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

func (sl *StringLiteral) Line() uint {
	return sl.Token.Line
}
