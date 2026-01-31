package ast

import (
	"fmt"

	"github.com/textwire/textwire/v3/token"
)

type UseStmt struct {
	BaseNode
	Name   *StringLiteral // The relative path to the layout like 'layouts/main'
	Layout *Program       // Pointer to the layout program
}

func NewUseStmt(tok token.Token) *UseStmt {
	return &UseStmt{
		BaseNode: NewBaseNode(tok),
	}
}

func (us *UseStmt) statementNode() {}

func (us *UseStmt) String() string {
	return fmt.Sprintf(`@use(%s)`, us.Name.String())
}

func (us *UseStmt) Stmts() []Statement {
	if us.Layout == nil {
		return []Statement{}
	}

	return us.Layout.Stmts()
}
