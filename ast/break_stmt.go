package ast

import (
	"github.com/textwire/textwire/v2/token"
)

type BreakStmt struct {
	Token token.Token // The '@break' token
	Pos   token.Position
}

func (bs *BreakStmt) statementNode() {
}

func (bs *BreakStmt) TokenLiteral() string {
	return bs.Token.Literal
}

func (bs *BreakStmt) String() string {
	return bs.Token.Literal
}

func (bs *BreakStmt) Line() uint {
	return bs.Token.DebugLine
}

func (bs *BreakStmt) Position() token.Position {
	return bs.Pos
}
