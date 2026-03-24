package ast

import (
	"fmt"
	"strings"

	"github.com/textwire/textwire/v4/pkg/token"
)

type GlobalFuncName string

type argRules struct {
	Min int
	Max int
}

const (
	defined    GlobalFuncName = "defined"
	hasValue   GlobalFuncName = "hasValue"
	formatDate GlobalFuncName = "formatDate"
)

var GlobalFunctions = map[GlobalFuncName]argRules{
	defined:    {Min: 1, Max: 999},
	hasValue:   {Min: 1, Max: 999},
	formatDate: {Min: 2, Max: 2},
}

type GlobalCallExpr struct {
	BaseNode
	Name      GlobalFuncName
	Arguments []Expression
}

func NewGlobalCallExpr(tok token.Token, name GlobalFuncName) *GlobalCallExpr {
	return &GlobalCallExpr{
		BaseNode: NewBaseNode(tok),
		Name:     name,
	}
}

func (*GlobalCallExpr) expressionNode() {}
func (*GlobalCallExpr) segmentNode()    {}

func (gce *GlobalCallExpr) String() string {
	var args strings.Builder
	args.Grow(len(gce.Arguments))

	for i := range gce.Arguments {
		args.WriteString(gce.Arguments[i].String())

		if i < len(gce.Arguments)-1 {
			args.WriteString(", ")
		}
	}

	return fmt.Sprintf("(%s(%s))", gce.Name, args.String())
}
