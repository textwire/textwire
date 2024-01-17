package ast

import (
	"fmt"

	"github.com/textwire/textwire/token"
)

type LayoutStatement struct {
	Token   token.Token
	Path    *StringLiteral
	Program *Program
}

func (ls *LayoutStatement) statementNode() {
}

func (ls *LayoutStatement) TokenLiteral() string {
	return ls.Token.Literal
}

func (ls *LayoutStatement) String() string {
	return fmt.Sprintf(`{{ layout "%s" }}`, ls.Path.String())
}
