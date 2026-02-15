package ast

import (
	"fmt"
	"strings"

	"github.com/textwire/textwire/v3/pkg/token"
)

type GlobalCallExp struct {
	BaseNode
	Function  *Identifier // Function being called
	Arguments []Expression
}

func NewGlobalCallExp(tok token.Token, function *Identifier) *GlobalCallExp {
	return &GlobalCallExp{
		BaseNode: NewBaseNode(tok),
		Function: function,
	}
}

func (gce *GlobalCallExp) expressionNode() {}

func (gce *GlobalCallExp) String() string {
	var args strings.Builder
	args.Grow(len(gce.Arguments))

	for i := range gce.Arguments {
		args.WriteString(gce.Arguments[i].String())

		if i < len(gce.Arguments)-1 {
			args.WriteString(", ")
		}
	}

	return fmt.Sprintf("(%s(%s))", gce.Function, args.String())
}
