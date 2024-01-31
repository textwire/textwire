package ast

import (
	"fmt"

	"github.com/textwire/textwire/token"
)

type IndexExpression struct {
	Token token.Token // The '[' token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode() {
}

func (ie *IndexExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *IndexExpression) String() string {
	return fmt.Sprintf("(%s[%s])", ie.Left, ie.Index.String())
}

func (ie *IndexExpression) Line() uint {
	return ie.Token.Line
}
