package ast

import "github.com/textwire/textwire/v2/token"

type ContinueStmt struct {
	Token token.Token // The '@continue' token
	Pos   Position
}

func (cs *ContinueStmt) statementNode() {
}

func (cs *ContinueStmt) TokenLiteral() string {
	return cs.Token.Literal
}

func (cs *ContinueStmt) String() string {
	return cs.Token.Literal
}

func (cs *ContinueStmt) Line() uint {
	return cs.Token.Line
}

func (cs *ContinueStmt) Position() Position {
	return cs.Pos
}
