package ast

import "github.com/textwire/textwire/token"

type HTMLStatement struct {
	Token token.Token
}

func (hs *HTMLStatement) statementNode() {
}

func (hs *HTMLStatement) TokenLiteral() string {
	return hs.Token.Literal
}

func (hs *HTMLStatement) String() string {
	return hs.Token.Literal
}
