package ast

import (
	"fmt"

	"github.com/textwire/textwire/v3/pkg/token"
)

type ContinueifStmt struct {
	BaseNode
	Condition Expression
}

func NewContinueIfStmt(tok token.Token) *ContinueifStmt {
	return &ContinueifStmt{
		BaseNode: NewBaseNode(tok),
	}
}

func (cis *ContinueifStmt) statementNode() {}

func (cis *ContinueifStmt) String() string {
	return fmt.Sprintf("%s(%s)", cis.Token.Literal, cis.Condition)
}
