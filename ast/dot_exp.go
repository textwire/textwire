package ast

import (
	"fmt"

	"github.com/textwire/textwire/v2/token"
)

type DotExp struct {
	Token token.Token // The dot token
	Left  Expression  // -->x.y
	Key   Expression  // x.y<--
	Pos   token.Position
}

func NewDotExp(tok token.Token, left Expression) *DotExp {
	return &DotExp{
		Token: tok, // "."
		Pos:   tok.Pos,
		Left:  left,
	}
}

func (de *DotExp) expressionNode() {}

func (de *DotExp) Tok() *token.Token {
	return &de.Token
}

func (de *DotExp) String() string {
	return fmt.Sprintf("(%s.%s)", de.Left, de.Key)
}

func (de *DotExp) Line() uint {
	return de.Token.ErrorLine()
}

func (de *DotExp) Position() token.Position {
	return de.Pos
}
