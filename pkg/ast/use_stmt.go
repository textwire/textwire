package ast

import (
	"fmt"

	"github.com/textwire/textwire/v3/pkg/token"
)

type UseStmt struct {
	BaseNode
	Name       *StringLiteral         // Relative path to the layout like 'layouts/main'
	LayoutProg *Program               // AST node of the layout file Name
	Inserts    map[string]*InsertStmt // @use connection to @insert directives
}

func NewUseStmt(tok token.Token) *UseStmt {
	return &UseStmt{
		BaseNode: NewBaseNode(tok),
	}
}

func (us *UseStmt) statementNode() {}

func (us *UseStmt) String() string {
	return fmt.Sprintf(`@use(%s)`, us.Name)
}

func (us *UseStmt) Stmts() []Statement {
	if us.LayoutProg == nil {
		return []Statement{}
	}

	return us.LayoutProg.Stmts()
}
