package ast

import "github.com/textwire/textwire/token"

type ExpressionStmt struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStmt) statementNode() {
}

func (es *ExpressionStmt) TokenLiteral() string {
	return es.Token.Literal
}

func (es *ExpressionStmt) String() string {
	if es.Expression == nil {
		return ""
	}

	return es.Expression.String()
}

func (es *ExpressionStmt) Line() uint {
	return es.Token.Line
}
