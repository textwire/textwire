package ast

import (
	"fmt"

	"github.com/textwire/textwire/v3/pkg/token"
)

type BreakifStmt struct {
	BaseNode
	Condition Expression
}

func NewBreakIfStmt(tok token.Token) *BreakifStmt {
	return &BreakifStmt{
		BaseNode: NewBaseNode(tok),
	}
}

func (bis *BreakifStmt) statementNode() {}

func (bis *BreakifStmt) String() string {
	return fmt.Sprintf("%s(%s)", bis.Token.Literal, bis.Condition)
}
