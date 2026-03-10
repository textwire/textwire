package ast

import (
	"github.com/textwire/textwire/v3/pkg/token"
)

type ChunkKind string

const (
	ChunkKindDirective ChunkKind = "directive"
	ChunkKindText      ChunkKind = "text"
	ChunkKindEmbedded  ChunkKind = "embedded"
	ChunkKindBlock     ChunkKind = "block"
)

// Chunk is a top lever node that your program is composed of.
// It can be a directive, text, or embedded code.
// - Embedded `{{ ... }}`.
// - Directive `@use('base')`.
// - Text `<h1>Hello</h1>`.
// - Block is a collection of chunks, 0 or more
type Chunk interface {
	Node
	chunkNode()
	Kind() ChunkKind
}

type Node interface {
	Tok() *token.Token
	SetTok(token.Token)
	String() string
	Line() uint
	Position() token.Position
	SetEndPosition(pos token.Position)
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
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
