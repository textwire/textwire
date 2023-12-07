package ast

import "github.com/textwire/textwire/token"

type HTMLStatement struct {
	Token token.Token
}

func (h *HTMLStatement) statementNode() {
}

func (h *HTMLStatement) TokenLiteral() string {
	return h.Token.Literal
}

func (h *HTMLStatement) String() string {
	return h.Token.Literal
}
