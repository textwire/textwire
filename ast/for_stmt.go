package ast

import (
	"bytes"

	"github.com/textwire/textwire/v2/token"
)

type ForStmt struct {
	BaseNode
	Init        Statement  // The initialization statement; or nil
	Condition   Expression // The condition expression; or nil
	Post        Statement  // The post iteration statement; or nil
	Alternative *BlockStmt // The @else block
	Block       *BlockStmt
}

func NewForStmt(tok token.Token) *ForStmt {
	return &ForStmt{
		BaseNode: NewBaseNode(tok),
	}
}

func (fs *ForStmt) statementNode() {}

func (fs *ForStmt) LoopBodyBlock() *BlockStmt {
	return fs.Block
}

func (fs *ForStmt) String() string {
	var out bytes.Buffer

	out.WriteString("@for(")
	out.WriteString(fs.Init.String() + "; ")
	out.WriteString(fs.Condition.String() + "; ")
	out.WriteString(fs.Post.String())
	out.WriteString(")\n")

	out.WriteString(fs.Block.String() + "\n")

	if fs.Alternative != nil {
		out.WriteString("@else\n")
		out.WriteString(fs.Alternative.String() + "\n")
	}

	out.WriteString("@end\n")

	return out.String()
}

func (fs *ForStmt) Stmts() []Statement {
	res := make([]Statement, 0)

	if fs.Block != nil {
		res = append(res, fs.Block.Stmts()...)
	}

	if fs.Alternative != nil {
		res = append(res, fs.Alternative.Stmts()...)
	}

	return res
}
