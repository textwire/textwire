package ast

import (
	"bytes"

	"github.com/textwire/textwire/v2/token"
)

type SlotStmt struct {
	Token token.Token    // The '@slot' token
	Name  *StringLiteral // when empty string literal, it means default slot
	Body  *BlockStmt     // optional block statement, can be nil
	Pos   token.Position
}

func (ss *SlotStmt) statementNode() {
}

func (ss *SlotStmt) Stmts() []Statement {
	return ss.Body.Statements
}

func (ss *SlotStmt) Tok() *token.Token {
	return &ss.Token
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
	return ss.Token.ErrorLine()
}

func (ss *SlotStmt) Position() token.Position {
	return ss.Pos
}
