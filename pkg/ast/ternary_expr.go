package ast

import (
	"fmt"

	"github.com/textwire/textwire/v3/pkg/token"
)

type TernaryExpr struct {
	BaseNode
	Cond     Expression
	IfExpr   Expression // true ? <IfExpr> : <ElseExpr>
	ElseExpr Expression // false ? <IfExpr> : <ElseExpr>
}

func NewTernaryExpr(tok token.Token, cond Expression) *TernaryExpr {
	return &TernaryExpr{
		BaseNode: NewBaseNode(tok),
		Cond:     cond,
	}
}

func (_ *TernaryExpr) expressionNode() {}

func (te *TernaryExpr) String() string {
	return fmt.Sprintf("(%s ? %s : %s)", te.Cond, te.IfExpr, te.ElseExpr)
}
