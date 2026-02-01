package ast

import (
	"fmt"
	"strings"

	"github.com/textwire/textwire/v3/token"
)

type CallExp struct {
	BaseNode
	Receiver  Expression  // Receiver of the call
	Function  *Identifier // Function being called
	Arguments []Expression
}

func NewCallExp(tok token.Token, receiver Expression, function *Identifier) *CallExp {
	return &CallExp{
		BaseNode: NewBaseNode(tok),
		Receiver: receiver,
		Function: function,
	}
}

func (ce *CallExp) expressionNode() {}

func (ce *CallExp) String() string {
	var args strings.Builder
	args.Grow(len(ce.Arguments) + (2 * len(ce.Arguments)))

	for i, arg := range ce.Arguments {
		args.WriteString(arg.String())

		if i < len(ce.Arguments)-1 {
			args.WriteString(", ")
		}
	}

	return fmt.Sprintf("(%s.%s(%s))", ce.Receiver, ce.Function, args.String())
}
