package ast

import (
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
	return ls.Program.String()
}
