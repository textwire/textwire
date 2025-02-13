package ast

import "github.com/textwire/textwire/v2/token"

type ExpressionStmt struct {
	Token      token.Token
	Expression Expression
	Pos        Position
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
	return es.Token.DebugLine
}

func (es *ExpressionStmt) Position() Position {
	return es.Pos
}
