package ast

import (
	"github.com/textwire/textwire/v4/pkg/position"
	"github.com/textwire/textwire/v4/pkg/token"
)

type Node interface {
	Tok() *token.Token
	SetTok(token.Token)
	Pos() *position.Pos
	SetEndPosition(pos *position.Pos)
	String() string
}

// Chunk is a top lever node that your program is composed of.
// It can be a directive, text, block, or embedded code.
type Chunk interface {
	Node
	chunkNode()
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// Segment represents a node that can appear inside {{ ... }} separated
// by a semicolon. This includes all expressions and statements.
type Segment interface {
	Node
	segmentNode()
}

type LoopDirective interface {
	LoopBlock() *Block
}

type NodeWithChunks interface {
	Node
	AllChunks() []Chunk
}

type SlotDirective interface {
	Node
	Name() *StrExpr
	IsDefault() bool
	SetIsDefault(bool)
	Block() *Block
	SetBlock(*Block)
}
