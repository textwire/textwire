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

func (*IncStmt) statementNode() {}
func (*IncStmt) segmentNode()   {}

func (is *IncStmt) String() string {
	return fmt.Sprintf("(%s++)", is.Left)
}
