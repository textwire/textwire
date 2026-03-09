package ast

import (
	"fmt"

	"github.com/textwire/textwire/v3/pkg/token"
)

type DecStmt struct {
	BaseNode
	Left Expression
}

func NewDecStmt(tok token.Token, left Expression) *DecStmt {
	return &DecStmt{
		BaseNode: NewBaseNode(tok),
		Left:     left,
	}
}

func (pe *DecStmt) statementNode() {}

func (pe *DecStmt) String() string {
	return fmt.Sprintf("(%s--)", pe.Left)
}
