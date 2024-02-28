package ast

import (
	"bytes"

	"github.com/textwire/textwire/token"
)

type ForStmt struct {
	Token       token.Token // The '@for' token
	Init        Statement   // The initialization statement; or nil
	Condition   Expression  // The condition expression; or nil
	Post        Statement   // The post iteration statement; or nil
	Alternative *BlockStmt  // The @else block
	Block       *BlockStmt
}

func (fs *ForStmt) statementNode() {
}

func (fs *ForStmt) TokenLiteral() string {
	return fs.Token.Literal
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

func (fs *ForStmt) Line() uint {
	return fs.Token.Line
}
