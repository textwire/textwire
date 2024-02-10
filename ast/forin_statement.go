package ast

import (
	"bytes"

	"github.com/textwire/textwire/token"
)

type ForinStatement struct {
	Token token.Token // The '@for' token
	Var   *Identifier // The variable name
	Array Expression  // The array expression
	Block *BlockStatement
}

func (fs *ForinStatement) statementNode() {
}

func (fs *ForinStatement) TokenLiteral() string {
	return fs.Token.Literal
}

func (fs *ForinStatement) String() string {
	var out bytes.Buffer

	out.WriteString("@for(" + fs.Var.String() + " in ")
	out.WriteString(fs.Array.String() + ")\n")
	out.WriteString(fs.Block.String() + "\n")

	out.WriteString("@end\n")

	return out.String()
}

func (fs *ForinStatement) Line() uint {
	return fs.Token.Line
}
