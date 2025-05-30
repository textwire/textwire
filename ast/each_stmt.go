package ast

import (
	"bytes"

	"github.com/textwire/textwire/v2/token"
)

type EachStmt struct {
	BaseNode
	Var         *Identifier // The variable name
	Array       Expression  // The array to loop over
	Alternative *BlockStmt  // The @else block
	Block       *BlockStmt
}

func NewEachStmt(tok token.Token) *EachStmt {
	return &EachStmt{
		BaseNode: NewBaseNode(tok),
	}
}

func (es *EachStmt) statementNode() {}

func (es *EachStmt) LoopBodyBlock() *BlockStmt {
	return es.Block
}

func (es *EachStmt) String() string {
	var out bytes.Buffer

	out.WriteString("@each(")
	out.WriteString(es.Var.String())
	out.WriteString(" in ")
	out.WriteString(es.Array.String())
	out.WriteString(")\n")
	out.WriteString(es.Block.String() + "\n")

	if es.Alternative != nil {
		out.WriteString("@else\n")
		out.WriteString(es.Alternative.String() + "\n")
	}

	out.WriteString("@end\n")

	return out.String()
}

func (es *EachStmt) Stmts() []Statement {
	res := make([]Statement, 0)

	if es.Block != nil {
		res = append(res, es.Block.Stmts()...)
	}

	if es.Alternative != nil {
		res = append(res, es.Alternative.Stmts()...)
	}

	return res
}
