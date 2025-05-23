package ast

import (
	"github.com/textwire/textwire/v2/token"
)

type BreakIfStmt struct {
	Token     token.Token // The '@breakIf' token
	Condition Expression
	Pos       token.Position
}

func NewBreakIfStmt(tok token.Token) *BreakIfStmt {
	return &BreakIfStmt{
		Token: tok, // "@breakIf"
		Pos:   tok.Pos,
	}
}

func (bis *BreakIfStmt) statementNode() {}

func (bis *BreakIfStmt) Tok() *token.Token {
	return &bis.Token
}

func (bis *BreakIfStmt) String() string {
	return bis.Token.Literal + "(" + bis.Condition.String() + ")"
}

func (bis *BreakIfStmt) Line() uint {
	return bis.Token.ErrorLine()
}

func (bis *BreakIfStmt) Position() token.Position {
	return bis.Pos
}
