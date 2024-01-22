package ast

import (
	"fmt"

	"github.com/textwire/textwire/token"
)

type IndexExpression struct {
	Token token.Token
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
