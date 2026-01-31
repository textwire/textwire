package ast

import "github.com/textwire/textwire/v3/token"

type ContinueIfStmt struct {
	BaseNode
	Condition Expression
}

func NewContinueIfStmt(tok token.Token) *ContinueIfStmt {
	return &ContinueIfStmt{
		BaseNode: NewBaseNode(tok),
	}
}

func (cis *ContinueIfStmt) statementNode() {}

func (cis *ContinueIfStmt) String() string {
	return cis.Token.Literal + "(" + cis.Condition.String() + ")"
}
