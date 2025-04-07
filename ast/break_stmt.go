package ast

import (
	"github.com/textwire/textwire/v2/token"
)

type BreakStmt struct {
	Token token.Token // The '@break' token
	Pos   token.Position
}

func (bs *BreakStmt) statementNode() {}

func (bs *BreakStmt) Tok() *token.Token {
	return &bs.Token
}

func (bs *BreakStmt) String() string {
	return bs.Token.Literal
}

func (bs *BreakStmt) Line() uint {
	return bs.Token.ErrorLine()
}

func (bs *BreakStmt) Position() token.Position {
	return bs.Pos
}
