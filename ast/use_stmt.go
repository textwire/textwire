package ast

import (
	"fmt"

	"github.com/textwire/textwire/v2/token"
)

type UseStmt struct {
	BaseNode
	Name    *StringLiteral // The relative path to the layout like 'layouts/main'
	Program *Program
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
