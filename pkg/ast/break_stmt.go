package ast

import (
	"github.com/textwire/textwire/v3/pkg/token"
)

type BreakStmt struct {
	BaseNode
}

func NewBreakStmt(tok token.Token) *BreakStmt {
	return &BreakStmt{
		BaseNode: NewBaseNode(tok),
	}
}

func (bs *BreakStmt) statementNode() {}

func (bs *BreakStmt) String() string {
	return bs.Token.Literal
}
