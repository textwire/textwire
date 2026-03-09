package ast

import "github.com/textwire/textwire/v3/pkg/token"

type ContinueDir struct {
	BaseNode
}

func NewContinueDir(tok token.Token) *ContinueDir {
	return &ContinueDir{
		BaseNode: NewBaseNode(tok),
	}
}

func (_ *ContinueDir) chunkNode() {}

func (_ *ContinueDir) Kind() ChunkKind {
	return ChunkKindDirective
}

func (cd *ContinueDir) String() string {
	return cd.Token.Lit
}
