package ast

import (
	"github.com/textwire/textwire/v2/token"
)

type AssignStmt struct {
	BaseNode
	Name  *Identifier
	Value Expression
}

func NewAssignStmt(tok token.Token, name *Identifier) *AssignStmt {
	return &AssignStmt{
		BaseNode: NewBaseNode(tok),
		Name:     name,
	}
}

func (as *AssignStmt) statementNode() {}

func (as *AssignStmt) String() string {
	return as.Name.String() + " = " + as.Value.String()
}
