package ast

import "github.com/textwire/textwire/token"

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {
}

func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

func (es *ExpressionStatement) String() string {
	if es.Expression == nil {
		return ""
	}

	return "{{ " + es.Expression.String() + " }}"
}

func (es *ExpressionStatement) LineNum() uint {
	return es.Token.Line
}
