package ast

import (
	"fmt"

	"github.com/textwire/textwire/v3/token"
)

type ReserveStmt struct {
	BaseNode
	Name *StringLiteral
	// Fallback is the second argument; nil if not present
	Fallback Expression
}

func NewReserveStmt(tok token.Token) *ReserveStmt {
	return &ReserveStmt{
		BaseNode: NewBaseNode(tok),
	}
}

func (rs *ReserveStmt) statementNode() {}

func (rs *ReserveStmt) String() string {
	if rs.Fallback == nil {
		return fmt.Sprintf(`@reserve("%s")`, rs.Name)
	}
	return fmt.Sprintf(`@reserve("%s", %s)`, rs.Name, rs.Fallback)
}
