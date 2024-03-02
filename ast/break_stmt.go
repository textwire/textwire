package ast

import (
	"github.com/textwire/textwire/token"
)

type BreakStmt struct {
	Token token.Token // The '@break' token
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
	return bs.Token.Line
}
