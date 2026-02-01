package ast

import (
	"fmt"
	"strings"

	"github.com/textwire/textwire/v3/token"
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
	var out strings.Builder
	out.Grow(26)

	fmt.Fprintf(&out, "@each(%s in %s)\n%s\n", es.Var, es.Array, es.Block)

	if es.Alternative != nil {
		out.WriteString("@else\n")
		out.WriteString(es.Alternative.String() + "\n")
	}

	out.WriteString("@end\n")

	return out.String()
}

func (es *EachStmt) Stmts() []Statement {
	stmts := make([]Statement, 0)

	if es.Block != nil {
		stmts = append(stmts, es.Block.Stmts()...)
	}

	if es.Alternative != nil {
		stmts = append(stmts, es.Alternative.Stmts()...)
	}

	return stmts
}
