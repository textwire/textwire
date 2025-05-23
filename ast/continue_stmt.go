package ast

import "github.com/textwire/textwire/v2/token"

type ContinueStmt struct {
	Token token.Token // The '@continue' token
	Pos   token.Position
}

func NewBreakStmt(tok token.Token) *BreakStmt {
	return &BreakStmt{
		Token: tok, // "@break"
		Pos:   tok.Pos,
	}
}

func NewContinueStmt(tok token.Token) *ContinueStmt {
	return &ContinueStmt{
		Token: tok, // "@continue"
		Pos:   tok.Pos,
	}
}

func (cs *ContinueStmt) statementNode() {}

func (cs *ContinueStmt) Tok() *token.Token {
	return &cs.Token
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
