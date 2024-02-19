package ast

import (
	"fmt"

	"github.com/textwire/textwire/token"
)

type DotExpression struct {
	Token    token.Token // The field token
	Receiver Expression  // -->x.y
	Key      Expression  // x.y<--
}

func (de *DotExpression) expressionNode() {
}

func (de *DotExpression) TokenLiteral() string {
	return de.Token.Literal
}

func (de *DotExpression) String() string {
	return fmt.Sprintf("(%s.%s)", de.Receiver, de.Key)
}

func (de *DotExpression) Line() uint {
	return de.Token.Line
}
