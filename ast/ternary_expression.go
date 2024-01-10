package ast

import (
	"fmt"

	"github.com/textwire/textwire/token"
)

type TernaryExpression struct {
	Token     token.Token
	Condition Expression
	TrueExp   Expression
	FalseExp  Expression
}

func (te *TernaryExpression) expressionNode() {
}

func (te *TernaryExpression) TokenLiteral() string {
	return te.Token.Literal
}

func (te *TernaryExpression) String() string {
	return fmt.Sprintf("(%s ? %s : %s)", te.Condition, te.TrueExp, te.FalseExp)
}
