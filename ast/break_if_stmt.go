package ast

import (
	"github.com/textwire/textwire/v3/token"
)

type BreakIfStmt struct {
	BaseNode
	Condition Expression
}

func NewBreakIfStmt(tok token.Token) *BreakIfStmt {
	return &BreakIfStmt{
		BaseNode: NewBaseNode(tok),
	}
}

func (bis *BreakIfStmt) statementNode() {}

func (bis *BreakIfStmt) String() string {
	return bis.Token.Literal + "(" + bis.Condition.String() + ")"
}
