package ast

import (
	"bytes"
	"fmt"

	"github.com/textwire/textwire/v3/token"
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
	var args bytes.Buffer

	for i, arg := range gce.Arguments {
		args.WriteString(arg.String())

		if i < len(gce.Arguments)-1 {
			args.WriteString(", ")
		}
	}

	return fmt.Sprintf("%s(%s))", gce.Function, args.String())
}
