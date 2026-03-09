package ast

import (
	"github.com/textwire/textwire/v3/pkg/token"
)

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Node interface {
	Tok() *token.Token
	String() string
	Line() uint
	Position() token.Position
	SetEndPosition(pos token.Position)
}

type LoopCommand interface {
	LoopBlock() *BlockStmt
}

type NodeWithStmts interface {
	Node
	Stmts() []Statement
}

type SlotCommand interface {
	Node
	Name() *StrLit
	IsDefault() bool
	SetIsDefault(bool)
	Block() *BlockStmt
	SetBlock(*BlockStmt)
}
