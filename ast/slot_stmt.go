package ast

import (
	"bytes"

	"github.com/textwire/textwire/token"
)

type SlotStmt struct {
	Token token.Token // The '@slot' token
	Name  *StringLiteral
	Body  *BlockStmt
}

func (ss *SlotStmt) statementNode() {
}

func (ss *SlotStmt) TokenLiteral() string {
	return ss.Token.Literal
}

func (ss *SlotStmt) String() string {
	var out bytes.Buffer

	out.WriteString("@slot(")
	out.WriteString(ss.Name.String())
	out.WriteString(")\n")

	out.WriteString(ss.Body.String())

	out.WriteString("\n@end")

	return out.String()
}

func (ss *SlotStmt) Line() uint {
	return ss.Token.Line
}
