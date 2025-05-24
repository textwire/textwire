package ast

import (
	"github.com/textwire/textwire/v2/token"
)

type Node interface {
	Tok() *token.Token
	String() string
	Line() uint
	Position() token.Position
	SetEndPosition(pos token.Position)
}

// TODO: See if you use it
// StatementsContainer helps to identify nodes that nest other statements.
type StatementsContainer interface {
	Stmts() []Statement
}

type LoopStmt interface {
	LoopBodyBlock() *BlockStmt
}
