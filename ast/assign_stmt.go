package ast

import (
	"github.com/textwire/textwire/v2/token"
)

type AssignStmt struct {
	BaseNode
	Left  *Identifier
	Right Expression
}

func NewAssignStmt(tok token.Token, left *Identifier) *AssignStmt {
	return &AssignStmt{
		BaseNode: NewBaseNode(tok),
		Left:     left,
	}
}

func (as *AssignStmt) statementNode() {}

func (as *AssignStmt) String() string {
	return as.Left.String() + " = " + as.Right.String()
}
