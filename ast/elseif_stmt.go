package ast

import (
	"fmt"
	"strings"

	"github.com/textwire/textwire/v3/token"
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
	var out strings.Builder

	fmt.Fprintf(&out, "@elseif(%s)\n", eis.Condition)
	out.WriteString(eis.Consequence.String())

	return out.String()
}

func (eis *ElseIfStmt) Stmts() []Statement {
	if eis.Consequence == nil {
		return []Statement{}
	}

	return eis.Consequence.Stmts()
}
