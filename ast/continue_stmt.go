package ast

import "github.com/textwire/textwire/v2/token"

type ContinueStmt struct {
	Token token.Token // The '@continue' token
	Pos   token.Position
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
	return cs.Token.ErrorLine()
}

func (cs *ContinueStmt) Position() token.Position {
	return cs.Pos
}
