package ast

import (
	"bytes"

	"github.com/textwire/textwire/token"
)

type ForStatement struct {
	Token     token.Token // The '@for' token
	Init      Statement   // The initialization statement; or nil
	Condition Expression  // The condition expression; or nil
	Post      Statement   // The post iteration statement; or nil
	Body      *BlockStatement
}

func (fs *ForStatement) statementNode() {
}

func (fs *ForStatement) TokenLiteral() string {
	return fs.Token.Literal
}

func (fs *ForStatement) String() string {
	var out bytes.Buffer

	out.WriteString("@for(")
	out.WriteString(fs.Init.String() + "; ")
	out.WriteString(fs.Condition.String() + "; ")
	out.WriteString(fs.Post.String())
	out.WriteString(")\n")

	out.WriteString(fs.Body.String() + "\n")

	out.WriteString("@end\n")

	return out.String()
}

func (fs *ForStatement) Line() uint {
	return fs.Token.Line
}
