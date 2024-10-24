package ast

import "github.com/textwire/textwire/v2/token"

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
