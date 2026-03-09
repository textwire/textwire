package ast

import (
	"fmt"

	"github.com/textwire/textwire/v3/pkg/token"
)

type IncStmt struct {
	BaseNode
	Left Expression
}

func NewIncStmt(tok token.Token, left Expression) *IncStmt {
	return &IncStmt{
		BaseNode: NewBaseNode(tok),
		Left:     left,
	}
}

func (pe *IncStmt) statementNode() {}

func (pe *IncStmt) String() string {
	return fmt.Sprintf("(%s++)", pe.Left)
}
