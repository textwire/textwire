package ast

import (
	"bytes"

	"github.com/textwire/textwire/token"
)

type ReturnStatement struct {
	Token token.Token
	Value Expression
}

func (rs *ReturnStatement) statementNode() {
}

func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.Value != nil {
		out.WriteString(rs.Value.String())
	}

	return out.String()
}
