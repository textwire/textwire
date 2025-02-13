package ast

import (
	"fmt"

	"github.com/textwire/textwire/v2/token"
)

type ReserveStmt struct {
	Token  token.Token // The '@reserve' token
	Insert *InsertStmt // The insert statement; nil if not yet parsed
	Name   *StringLiteral
	Pos    Position
}

func (rs *ReserveStmt) statementNode() {
}

func (rs *ReserveStmt) TokenLiteral() string {
	return rs.Token.Literal
}

func (rs *ReserveStmt) String() string {
	return fmt.Sprintf(`@reserve("%s")`, rs.Name.String())
}

func (rs *ReserveStmt) Line() uint {
	return rs.Token.DebugLine
}

func (rs *ReserveStmt) Position() Position {
	return rs.Pos
}
