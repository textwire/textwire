package ast

import (
	"fmt"

	"github.com/textwire/textwire/v2/token"
)

type PrefixExp struct {
	Token    token.Token // The '!' or '-' token
	Operator string
	Right    Expression
	Pos      token.Position
}

func (pe *PrefixExp) expressionNode() {
}

func (pe *PrefixExp) TokenLiteral() string {
	return pe.Token.Literal
}

func (pe *PrefixExp) String() string {
	return fmt.Sprintf("(%s%s)", pe.Operator, pe.Right)
}

func (pe *PrefixExp) Line() uint {
	return pe.Token.ErrorLine()
}

func (pe *PrefixExp) Position() token.Position {
	return pe.Pos
}
