package ast

import "github.com/textwire/textwire/v2/token"

type ExpressionStmt struct {
	Token      token.Token
	Expression Expression
	Pos        token.Position
}

func NewExpressionStmt(tok token.Token, exp Expression) *ExpressionStmt {
	return &ExpressionStmt{
		Token:      tok,
		Pos:        tok.Pos,
		Expression: exp,
	}
}

func (es *ExpressionStmt) statementNode() {}

func (es *ExpressionStmt) Tok() *token.Token {
	return &es.Token
}

func (es *ExpressionStmt) String() string {
	if es.Expression == nil {
		return ""
	}

	return es.Expression.String()
}

func (es *ExpressionStmt) Line() uint {
	return es.Token.ErrorLine()
}

func (es *ExpressionStmt) Position() token.Position {
	return es.Pos
}
