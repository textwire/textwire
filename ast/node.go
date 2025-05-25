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

type LoopStmt interface {
	LoopBodyBlock() *BlockStmt
}
