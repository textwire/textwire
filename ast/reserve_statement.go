package ast

import (
	"fmt"

	"github.com/textwire/textwire/token"
)

type ReserveStatement struct {
	Token  token.Token // The '@reserve' token
	Insert *InsertStmt // The insert statement; nil if not yet parsed
	Name   *StringLiteral
}

func (rs *ReserveStatement) statementNode() {
}

func (rs *ReserveStatement) TokenLiteral() string {
	return rs.Token.Literal
}

func (rs *ReserveStatement) String() string {
	return fmt.Sprintf(`@reserve("%s")`, rs.Name.String())
}

func (rs *ReserveStatement) Line() uint {
	return rs.Token.Line
}
