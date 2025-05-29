package ast

import (
	"bytes"

	"github.com/textwire/textwire/v2/token"
)

type IfStmt struct {
	BaseNode
	Condition    Expression    // The truthy condition
	Consequence  *BlockStmt    // The 'then' block
	Alternative  *BlockStmt    // The @else block
	Alternatives []*ElseIfStmt // The @elseif blocks
}

func NewIfStmt(tok token.Token) *IfStmt {
	return &IfStmt{
		BaseNode: NewBaseNode(tok),
	}
}

func (is *IfStmt) statementNode() {}

func (is *IfStmt) String() string {
	var out bytes.Buffer

	out.WriteString("@if(" + is.Condition.String() + ")\n")

	out.WriteString(is.Consequence.String())

	for _, e := range is.Alternatives {
		out.WriteString(e.String())
	}

	if is.Alternative != nil {
		out.WriteString("@else\n")
		out.WriteString(is.Alternative.String() + "\n")
	}

	out.WriteString("@end\n")

	return out.String()
}

func (is *IfStmt) Stmts() []Statement {
	res := make([]Statement, 0)

	if is.Consequence != nil {
		res = append(res, is.Consequence.Stmts()...)
	}

	if is.Alternative != nil {
		res = append(res, is.Alternative.Stmts()...)
	}

	for _, e := range is.Alternatives {
		res = append(res, e.Stmts()...)
	}

	return res
}
