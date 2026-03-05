package ast

import (
	"fmt"

	"github.com/textwire/textwire/v3/pkg/token"
)

type AssignStmt struct {
	BaseNode
	Left  Expression
	Right Expression
}

func NewAssignStmt(tok token.Token, left Expression) *AssignStmt {
	return &AssignStmt{
		BaseNode: NewBaseNode(tok),
		Left:     left,
	}
}

func (as *AssignStmt) statementNode() {}

func (as *AssignStmt) String() string {
	return fmt.Sprintf("%s = %s", as.Left, as.Right)
}
