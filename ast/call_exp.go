package ast

import (
	"bytes"
	"fmt"

	"github.com/textwire/textwire/token"
)

type CallExp struct {
	Token     token.Token // Function identifier token
	Receiver  Expression  // Receiver of the call
	Function  *Identifier // Function being called
	Arguments []Expression
}

func (ce *CallExp) expressionNode() {
}

func (ce *CallExp) TokenLiteral() string {
	return ce.Token.Literal
}

func (ce *CallExp) String() string {
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

func (ce *CallExp) Line() uint {
	return ce.Token.Line
}
