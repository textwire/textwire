package ast

import (
	"fmt"

	"github.com/textwire/textwire/v2/token"
)

type ReserveStmt struct {
	Token  token.Token // The '@reserve' token
	Insert *InsertStmt // The insert statement; nil if not yet parsed
	Name   *StringLiteral
	Pos    token.Position
}

func NewReserveStmt(tok token.Token) *ReserveStmt {
	return &ReserveStmt{
		Token: tok, // "@reserve"
		Pos:   tok.Pos,
	}
}

func (rs *ReserveStmt) statementNode() {}

func (rs *ReserveStmt) Tok() *token.Token {
	return &rs.Token
}

func (rs *ReserveStmt) String() string {
	return fmt.Sprintf(`@reserve("%s")`, rs.Name.String())
}

func (rs *ReserveStmt) Line() uint {
	return rs.Token.ErrorLine()
}

func (rs *ReserveStmt) Position() token.Position {
	return rs.Pos
}
