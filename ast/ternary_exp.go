package ast

import (
	"fmt"

	"github.com/textwire/textwire/v3/token"
)

type TernaryExp struct {
	BaseNode
	Condition Expression
	IfBlock   Expression // true ? <IfBlock> : <ElseBlock>
	ElseBlock Expression // false ? <IfBlock> : <ElseBlock>
}

func NewTernaryExp(tok token.Token, cond Expression) *TernaryExp {
	return &TernaryExp{
		BaseNode:  NewBaseNode(tok),
		Condition: cond,
	}
}

func (te *TernaryExp) expressionNode() {}

func (te *TernaryExp) String() string {
	return fmt.Sprintf("(%s ? %s : %s)", te.Condition, te.IfBlock, te.ElseBlock)
}
