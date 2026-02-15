package ast

import (
	"fmt"
	"strings"

	"github.com/textwire/textwire/v3/pkg/token"
)

type EachStmt struct {
	BaseNode
	Var       *Identifier // Variable name
	Array     Expression  // Array to loop over
	ElseBlock *BlockStmt  // @else<ElseBlock>@end
	Block     *BlockStmt
}

func NewEachStmt(tok token.Token) *EachStmt {
	return &EachStmt{
		BaseNode: NewBaseNode(tok),
	}
}

func (es *EachStmt) statementNode() {}

func (es *EachStmt) LoopBlock() *BlockStmt {
	if es.Block == nil {
		panic("Block must not be nil on EachStmt when calling LoopBlock()")
	}
	return es.Block
}

func (es *EachStmt) String() string {
	var out strings.Builder
	out.Grow(26)

	fmt.Fprintf(&out, "@each(%s in %s)\n%s\n", es.Var, es.Array, es.Block)

	if es.ElseBlock != nil {
		out.WriteString("@else\n")
		out.WriteString(es.ElseBlock.String() + "\n")
	}

	out.WriteString("@end\n")

	return out.String()
}

func (es *EachStmt) Stmts() []Statement {
	stmts := make([]Statement, 0)
	if es.Block != nil {
		stmts = append(stmts, es.Block.Stmts()...)
	}

	if es.ElseBlock != nil {
		stmts = append(stmts, es.ElseBlock.Stmts()...)
	}

	return stmts
}
