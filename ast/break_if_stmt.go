package ast

import (
	"fmt"

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
	return fmt.Sprintf("%s(%s)", bis.Token.Literal, bis.Condition)
}
