package ast

import "github.com/textwire/textwire/v3/token"

type ContinueStmt struct {
	BaseNode
}

func NewContinueStmt(tok token.Token) *ContinueStmt {
	return &ContinueStmt{
		BaseNode: NewBaseNode(tok),
	}
}

func (cs *ContinueStmt) statementNode() {}

func (cs *ContinueStmt) String() string {
	return cs.Token.Literal
}
