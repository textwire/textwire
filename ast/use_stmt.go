package ast

import (
	"fmt"

	"github.com/textwire/textwire/v2/token"
)

type UseStmt struct {
	Token   token.Token    // The '@use' token
	Name    *StringLiteral // The relative path to the layout like 'layouts/main'
	Program *Program
	Pos     Position
}

func (us *UseStmt) statementNode() {
}

func (us *UseStmt) TokenLiteral() string {
	return us.Token.Literal
}

func (us *UseStmt) String() string {
	return fmt.Sprintf(`@use(%s)`, us.Name.String())
}

func (us *UseStmt) Line() uint {
	return us.Token.DebugLine
}

func (us *UseStmt) Position() Position {
	return us.Pos
}
