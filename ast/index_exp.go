package ast

import (
	"fmt"

	"github.com/textwire/textwire/token"
)

type IndexExp struct {
	Token token.Token // The '[' token
	Left  Expression
	Index Expression
}

func (ie *IndexExp) expressionNode() {
}

func (ie *IndexExp) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *IndexExp) String() string {
	return fmt.Sprintf("(%s[%s])", ie.Left, ie.Index.String())
}

func (ie *IndexExp) Line() uint {
	return ie.Token.Line
}
