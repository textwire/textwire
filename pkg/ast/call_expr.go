package ast

import (
	"fmt"
	"strings"

	"github.com/textwire/textwire/v3/pkg/token"
)

type CallExpr struct {
	BaseNode
	Receiver  Expression // Receiver of the call
	Function  *IdentExpr // Function being called
	Arguments []Expression
}

func NewCallExpr(tok token.Token, receiver Expression, function *IdentExpr) *CallExpr {
	return &CallExpr{
		BaseNode: NewBaseNode(tok),
		Receiver: receiver,
		Function: function,
	}
}

func (*CallExpr) expressionNode() {}

func (ce *CallExpr) String() string {
	var args strings.Builder
	args.Grow(len(ce.Arguments) + (2 * len(ce.Arguments)))

	for i := range ce.Arguments {
		args.WriteString(ce.Arguments[i].String())

		if i < len(ce.Arguments)-1 {
			args.WriteString(", ")
		}
	}

	return fmt.Sprintf("(%s.%s(%s))", ce.Receiver, ce.Function, args.String())
}
