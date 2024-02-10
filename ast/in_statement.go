package ast

import (
	"fmt"

	"github.com/textwire/textwire/token"
)

type InStatement struct {
	Token token.Token // The '@for' token
	Var   *Identifier // The variable name
	Array Expression  // The array to loop over
}

func (is *InStatement) statementNode() {
}

func (is *InStatement) TokenLiteral() string {
	return is.Token.Literal
}

func (is *InStatement) String() string {
	return fmt.Sprintf("%s %s in %s", is.TokenLiteral(),
		is.Var.String(), is.Array.String())
}

func (is *InStatement) Line() uint {
	return is.Token.Line
}
