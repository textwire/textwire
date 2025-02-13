package ast

import (
	"fmt"

	"github.com/textwire/textwire/v2/token"
)

type TernaryExp struct {
	Token       token.Token // The '?' token
	Condition   Expression
	Consequence Expression
	Alternative Expression
	Pos         Position
}

func (te *TernaryExp) expressionNode() {
}

func (te *TernaryExp) TokenLiteral() string {
	return te.Token.Literal
}

func (te *TernaryExp) String() string {
	return fmt.Sprintf("(%s ? %s : %s)", te.Condition, te.Condition, te.Alternative)
}

func (te *TernaryExp) Line() uint {
	return te.Token.DebugLine
}

func (te *TernaryExp) Position() Position {
	return te.Pos
}
