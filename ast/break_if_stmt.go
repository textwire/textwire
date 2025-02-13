package ast

import (
	"github.com/textwire/textwire/v2/token"
)

type BreakIfStmt struct {
	Token     token.Token // The '@breakIf' token
	Condition Expression
	Pos       Position
}

func (bis *BreakIfStmt) statementNode() {
}

func (bis *BreakIfStmt) TokenLiteral() string {
	return bis.Token.Literal
}

func (bis *BreakIfStmt) String() string {
	return bis.Token.Literal + "(" + bis.Condition.String() + ")"
}

func (bis *BreakIfStmt) Line() uint {
	return bis.Token.DebugLine
}

func (bis *BreakIfStmt) Position() Position {
	return bis.Pos
}
