package ast

import (
	"bytes"

	"github.com/textwire/textwire/v2/token"
)

type ElseIfStmt struct {
	BaseNode
	Condition   Expression
	Consequence *BlockStmt
}

func NewElseIfStmt(tok token.Token) *ElseIfStmt {
	return &ElseIfStmt{
		BaseNode: NewBaseNode(tok),
	}
}

func (eis *ElseIfStmt) statementNode() {}

func (eis *ElseIfStmt) String() string {
	var out bytes.Buffer

	out.WriteString("@elseif(" + eis.Condition.String() + ")\n")
	out.WriteString(eis.Consequence.String())

	return out.String()
}

func (eis *ElseIfStmt) Stmts() []Statement {
	if eis.Consequence == nil {
		return []Statement{}
	}

	return eis.Consequence.Stmts()
}
