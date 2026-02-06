package ast

import (
	"fmt"
	"strings"

	"github.com/textwire/textwire/v3/token"
)

type ElseIfStmt struct {
	BaseNode
	Condition Expression
	Block     *BlockStmt // @elseif()<Block>@end
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
	out.WriteString(eis.Block.String())

	return out.String()
}

func (eis *ElseIfStmt) Stmts() []Statement {
	if eis.Block == nil {
		return []Statement{}
	}

	return eis.Block.Stmts()
}
