package ast

import (
	"fmt"

	"github.com/textwire/textwire/v3/token"
)

type ReserveStmt struct {
	BaseNode
	Insert *InsertStmt // Insert statement; nil if not yet parsed
	Name   *StringLiteral
}

func NewReserveStmt(tok token.Token) *ReserveStmt {
	return &ReserveStmt{
		BaseNode: NewBaseNode(tok),
	}
}

func (rs *ReserveStmt) statementNode() {}

func (rs *ReserveStmt) String() string {
	return fmt.Sprintf(`@reserve("%s")`, rs.Name)
}
