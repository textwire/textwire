package ast

import (
	"github.com/textwire/textwire/v3/token"
)

type Node interface {
	Tok() *token.Token
	String() string
	Line() uint
	Position() token.Position
	SetEndPosition(pos token.Position)
}

type LoopStmt interface {
	LoopBodyBlock() *BlockStmt
}

type NodeWithStatements interface {
	Node
	Stmts() []Statement
}
