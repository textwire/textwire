package ast

import "github.com/textwire/textwire/v3/pkg/token"

// Text holds literal string of text. The token literal is its value.
type Text struct {
	BaseNode
}

func NewText(tok token.Token) *Text {
	return &Text{
		BaseNode: NewBaseNode(tok),
	}
}

func (_ *Text) chunkNode() {}

func (t *Text) String() string {
	return t.Token.Lit
}
