package ast

import (
	"bytes"

	"github.com/textwire/textwire/token"
)

type ElseIfStatement struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
}

func (eis *ElseIfStatement) statementNode() {
}

func (eis *ElseIfStatement) TokenLiteral() string {
	return eis.Token.Literal
}

func (eis *ElseIfStatement) String() string {
	var result bytes.Buffer

	result.WriteString("{{ else if " + eis.Condition.String() + " }}\n")
	result.WriteString(eis.Consequence.String())
	result.WriteString("\n{{ end }}\n")

	return result.String()
}
