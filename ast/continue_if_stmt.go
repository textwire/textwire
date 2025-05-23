package ast

import "github.com/textwire/textwire/v2/token"

type ContinueIfStmt struct {
	Token     token.Token // The '@continueIf' token
	Condition Expression
	Pos       token.Position
}

func NewContinueIfStmt(tok token.Token) *ContinueIfStmt {
	return &ContinueIfStmt{
		Token: tok, // "@continueIf"
		Pos:   tok.Pos,
	}
}

func (cis *ContinueIfStmt) statementNode() {}

func (cis *ContinueIfStmt) Tok() *token.Token {
	return &cis.Token
}

func (cis *ContinueIfStmt) String() string {
	return cis.Token.Literal + "(" + cis.Condition.String() + ")"
}

func (cis *ContinueIfStmt) Line() uint {
	return cis.Token.ErrorLine()
}

func (cis *ContinueIfStmt) Position() token.Position {
	return cis.Pos
}
