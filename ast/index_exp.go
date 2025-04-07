package ast

import (
	"fmt"

	"github.com/textwire/textwire/v2/token"
)

type IndexExp struct {
	Token token.Token // The '[' token
	Left  Expression
	Index Expression
	Pos   token.Position
}

func (ie *IndexExp) expressionNode() {}

func (ie *IndexExp) Tok() *token.Token {
	return &ie.Token
}

func (ie *IndexExp) String() string {
	return fmt.Sprintf("(%s[%s])", ie.Left, ie.Index.String())
}

func (ie *IndexExp) Line() uint {
	return ie.Token.ErrorLine()
}

func (ie *IndexExp) Position() token.Position {
	return ie.Pos
}
