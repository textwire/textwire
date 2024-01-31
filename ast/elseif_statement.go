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
	var out bytes.Buffer

	out.WriteString("@elseif(" + eis.Condition.String() + ")\n")
	out.WriteString(eis.Consequence.String())

	return out.String()
}

func (eis *ElseIfStatement) Line() uint {
	return eis.Token.Line
}
