package ast

import (
	"fmt"
	"strings"

	"github.com/textwire/textwire/v3/pkg/token"
)

type ForStmt struct {
	BaseNode
	Init      Statement  // Initialization statement; or nil
	Condition Expression // Condition expression; or nil
	Post      Statement  // Post iteration statement; or nil
	ElseBlock *BlockStmt // @else block
	Block     *BlockStmt
}

func NewForStmt(tok token.Token) *ForStmt {
	return &ForStmt{
		BaseNode: NewBaseNode(tok),
	}
}

func (fs *ForStmt) statementNode() {}

func (fs *ForStmt) LoopBlock() *BlockStmt {
	if fs.Block == nil {
		panic("Block must not be nil on ForStmt when calling LoopBlock()")
	}
	return fs.Block
}

func (fs *ForStmt) String() string {
	var out strings.Builder
	out.Grow(20)

	fmt.Fprintf(&out, "@for(%s; %s; %s)\n", fs.Init, fs.Condition, fs.Post)

	out.WriteString(fs.Block.String() + "\n")

	if fs.ElseBlock != nil {
		out.WriteString("@else\n")
		out.WriteString(fs.ElseBlock.String() + "\n")
	}

	out.WriteString("@end\n")

	return out.String()
}

func (fs *ForStmt) Stmts() []Statement {
	stmts := make([]Statement, 0)
	if fs.Block != nil {
		stmts = append(stmts, fs.Block.Stmts()...)
	}

	if fs.ElseBlock != nil {
		stmts = append(stmts, fs.ElseBlock.Stmts()...)
	}

	return stmts
}
