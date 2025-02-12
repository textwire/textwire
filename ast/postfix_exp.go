package ast

import (
	"fmt"

	"github.com/textwire/textwire/v2/token"
)

type PostfixExp struct {
	Token    token.Token // The '++' or '--' token
	Operator string
	Left     Expression
	Pos      Position
}

func (pe *PostfixExp) expressionNode() {
}

func (pe *PostfixExp) TokenLiteral() string {
	return pe.Token.Literal
}

func (pe *PostfixExp) String() string {
	return fmt.Sprintf("(%s%s)", pe.Left, pe.Operator)
}

func (pe *PostfixExp) Line() uint {
	return pe.Token.Line
}

func (pe *PostfixExp) Position() Position {
	return pe.Pos
}
