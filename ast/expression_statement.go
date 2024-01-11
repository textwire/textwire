package ast

import "github.com/textwire/textwire/token"

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (e *ExpressionStatement) statementNode() {
}

func (e *ExpressionStatement) TokenLiteral() string {
	return e.Token.Literal
}

func (e *ExpressionStatement) String() string {
	if e.Expression == nil {
		return ""
	}

	return "{{ " + e.Expression.String() + " }}"
}
