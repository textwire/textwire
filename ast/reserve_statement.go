package ast

import (
	"github.com/textwire/textwire/token"
)

type ReserveStatement struct {
	Token   token.Token
	Name    *StringLiteral
	Program *Program
}

func (rs *ReserveStatement) statementNode() {
}

func (rs *ReserveStatement) TokenLiteral() string {
	return rs.Token.Literal
}

func (rs *ReserveStatement) String() string {
	return rs.Program.String()
}
