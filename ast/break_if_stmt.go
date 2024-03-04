package ast

import (
	"github.com/textwire/textwire/token"
)

type BreakIfStmt struct {
	Token     token.Token // The '@breakIf' token
	Condition Expression
}

func (bis *BreakIfStmt) statementNode() {
}

func (bis *BreakIfStmt) TokenLiteral() string {
	return bis.Token.Literal
}

func (bis *BreakIfStmt) String() string {
	return bis.Token.Literal
}

func (bis *BreakIfStmt) Line() uint {
	return bis.Token.Line
}
