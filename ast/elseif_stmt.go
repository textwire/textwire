package ast

import (
	"bytes"

	"github.com/textwire/textwire/v2/token"
)

type ElseIfStmt struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStmt
	Pos         Position
}

func (eis *ElseIfStmt) statementNode() {
}

func (eis *ElseIfStmt) TokenLiteral() string {
	return eis.Token.Literal
}

func (eis *ElseIfStmt) String() string {
	var out bytes.Buffer

	out.WriteString("@elseif(" + eis.Condition.String() + ")\n")
	out.WriteString(eis.Consequence.String())

	return out.String()
}

func (eis *ElseIfStmt) Line() uint {
	return eis.Token.DebugLine
}

func (eis *ElseIfStmt) Position() Position {
	return eis.Pos
}
