package ast

import (
	"bytes"

	"github.com/textwire/textwire/v2/token"
)

type IfStmt struct {
	Token        token.Token   // The '@if' token
	Condition    Expression    // The truthy condition
	Consequence  *BlockStmt    // The 'then' block
	Alternative  *BlockStmt    // The @else block
	Alternatives []*ElseIfStmt // The @elseif blocks
	Pos          token.Position
}

func (is *IfStmt) statementNode() {
}

func (is *IfStmt) Tok() *token.Token {
	return &is.Token
}

func (is *IfStmt) String() string {
	var out bytes.Buffer

	out.WriteString("@if(" + is.Condition.String() + ")\n")

	out.WriteString(is.Consequence.String())

	for _, e := range is.Alternatives {
		out.WriteString(e.String())
	}

	if is.Alternative != nil {
		out.WriteString("@else\n")
		out.WriteString(is.Alternative.String() + "\n")
	}

	out.WriteString("@end\n")

	return out.String()
}

func (is *IfStmt) Line() uint {
	return is.Token.ErrorLine()
}

func (is *IfStmt) Position() token.Position {
	return is.Pos
}
