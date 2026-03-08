package ast

import (
	"fmt"

	"github.com/textwire/textwire/v3/pkg/token"
)

type StrLit struct {
	BaseNode
	Val string
}

func NewStrLit(tok token.Token, val string) *StrLit {
	return &StrLit{
		BaseNode: NewBaseNode(tok),
		Val:      val,
	}
}

func (sl *StrLit) expressionNode() {}

func (sl *StrLit) String() string {
	return fmt.Sprintf(`"%s"`, sl.Token.Lit)
}
