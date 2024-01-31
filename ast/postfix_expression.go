package ast

import (
	"fmt"

	"github.com/textwire/textwire/token"
)

type PostfixExpression struct {
	Token    token.Token // The '++' or '--' token
	Operator string
	Left     Expression
}

func (pe *PostfixExpression) expressionNode() {
}

func (pe *PostfixExpression) TokenLiteral() string {
	return pe.Token.Literal
}

func (pe *PostfixExpression) String() string {
	return fmt.Sprintf("(%s%s)", pe.Left, pe.Operator)
}

func (pe *PostfixExpression) Line() uint {
	return pe.Token.Line
}
