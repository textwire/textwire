package ast

import (
	"bytes"
	"fmt"

	"github.com/textwire/textwire/token"
)

type CallExpression struct {
	Token     token.Token // The receiver token
	Receiver  Expression  // The receiver of the call
	Function  *Identifier // The function being called
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {
}

func (ce *CallExpression) TokenLiteral() string {
	return ce.Token.Literal
}

func (ce *CallExpression) String() string {
	var args bytes.Buffer

	for i, arg := range ce.Arguments {
		args.WriteString(arg.String())

		if i < len(ce.Arguments)-1 {
			args.WriteString(", ")
		}
	}

	return fmt.Sprintf("%s.%s(%s)", ce.Receiver.String(),
		ce.Function.String(), args.String())
}

func (ce *CallExpression) Line() uint {
	return ce.Token.Line
}
