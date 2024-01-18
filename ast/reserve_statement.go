package ast

import (
	"fmt"

	"github.com/textwire/textwire/token"
)

type ReserveStatement struct {
	Token  token.Token
	Name   *StringLiteral
	Insert *InsertStatement
}

func (rs *ReserveStatement) statementNode() {
}

func (rs *ReserveStatement) TokenLiteral() string {
	return rs.Token.Literal
}

func (rs *ReserveStatement) String() string {
	return fmt.Sprintf(`{{ reserve "%s" }}`, rs.Name.String())
}
