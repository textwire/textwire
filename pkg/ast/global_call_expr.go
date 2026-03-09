package ast

import (
	"fmt"
	"strings"

	"github.com/textwire/textwire/v3/pkg/token"
)

type GlobalCallExpr struct {
	BaseNode
	Function  *IdentExpr // Function being called
	Arguments []Expression
}

func NewGlobalCallExpr(tok token.Token, function *IdentExpr) *GlobalCallExpr {
	return &GlobalCallExpr{
		BaseNode: NewBaseNode(tok),
		Function: function,
	}
}

func (_ *GlobalCallExpr) expressionNode() {}

func (gce *GlobalCallExpr) String() string {
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
