package ast

import (
	"bytes"

	"github.com/textwire/textwire/v2/token"
)

type SlotStmt struct {
	Token token.Token    // The '@slot' token
	Name  *StringLiteral // when empty string literal, it means default slot
	Body  *BlockStmt     // optional block statement, can be nil
	Pos   Position
}

func (ss *SlotStmt) statementNode() {
}

func (ss *SlotStmt) TokenLiteral() string {
	return ss.Token.Literal
}

func (ss *SlotStmt) String() string {
	var out bytes.Buffer

	if ss.Name.Value == "" {
		out.WriteString("@slot")
	} else {
		out.WriteString("@slot(")
		out.WriteString(ss.Name.String())
		out.WriteString(")")
	}

	if ss.Body != nil {
		out.WriteString("\n")
		out.WriteString(ss.Body.String())
		out.WriteString("\n@end")
	}

	return out.String()
}

func (ss *SlotStmt) Line() uint {
	return ss.Token.DebugLine
}

func (ss *SlotStmt) Position() Position {
	return ss.Pos
}
