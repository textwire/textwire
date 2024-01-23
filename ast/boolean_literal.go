package ast

import "github.com/textwire/textwire/token"

type BooleanLiteral struct {
	Token token.Token // The 'true' or 'false' token
	Value bool
}

func (bl *BooleanLiteral) expressionNode() {
}

func (bl *BooleanLiteral) TokenLiteral() string {
	return bl.Token.Literal
}

func (bl *BooleanLiteral) String() string {
	return bl.Token.Literal
}
