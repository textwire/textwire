package ast

import (
	"fmt"

	"github.com/textwire/textwire/v2/token"
)

type UseStmt struct {
	Token   token.Token    // The '@use' token
	Name    *StringLiteral // The relative path to the layout like 'layouts/main'
	Program *Program
	Pos     token.Position
}

func (us *UseStmt) statementNode() {}

func (us *UseStmt) Tok() *token.Token {
	return &us.Token
}

func (us *UseStmt) String() string {
	return fmt.Sprintf(`@use(%s)`, us.Name.String())
}

func (us *UseStmt) Line() uint {
	return us.Token.ErrorLine()
}

func (us *UseStmt) Position() token.Position {
	return us.Pos
}
