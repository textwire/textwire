package ast

import "github.com/textwire/textwire/token"

type LayoutStatement struct {
	Token token.Token
	Name  *StringLiteral
}

func (ls *LayoutStatement) statementNode() {
}

func (ls *LayoutStatement) TokenLiteral() string {
	return ls.Token.Literal
}

func (ls *LayoutStatement) String() string {
	return ls.Token.Literal
}
