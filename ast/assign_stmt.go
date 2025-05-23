package ast

import (
	"github.com/textwire/textwire/v2/token"
)

type AssignStmt struct {
	Token token.Token // The 'var' or identifier token
	Name  *Identifier
	Value Expression
	Pos   token.Position
}

func NewAssignStmt(tok token.Token, name *Identifier) *AssignStmt {
	return &AssignStmt{
		Token: tok, // identifier
		Pos:   tok.Pos,
		Name:  name,
	}
}

func (as *AssignStmt) statementNode() {}

func (as *AssignStmt) Tok() *token.Token {
	return &as.Token
}

func (as *AssignStmt) String() string {
	return as.Name.String() + " = " + as.Value.String()
}

func (as *AssignStmt) Line() uint {
	return as.Token.ErrorLine()
}

func (as *AssignStmt) Position() token.Position {
	return as.Pos
}
