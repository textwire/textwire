package ast

import "github.com/textwire/textwire/v3/pkg/token"

type ExpressionStmt struct {
	BaseNode
	Expression Expression
}

func NewExpressionStmt(tok token.Token, exp Expression) *ExpressionStmt {
	return &ExpressionStmt{
		BaseNode:   NewBaseNode(tok),
		Expression: exp,
	}
}

func (es *ExpressionStmt) statementNode() {}

func (es *ExpressionStmt) String() string {
	if es.Expression == nil {
		return ""
	}

	return es.Expression.String()
}
