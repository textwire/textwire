package ast

import (
	"github.com/textwire/textwire/v3/pkg/position"
	"github.com/textwire/textwire/v3/pkg/token"
)

// BaseNode is the main base for all the AST nodes.
type BaseNode struct {
	Token token.Token
	pos   *position.Pos
}

func NewBaseNode(tok token.Token) BaseNode {
	return BaseNode{
		Token: tok,
		pos:   tok.Pos,
	}
}

func (bn *BaseNode) Tok() *token.Token {
	return &bn.Token
}

func (bn *BaseNode) Pos() *position.Pos {
	return bn.pos
}

func (bn *BaseNode) SetEndPosition(pos *position.Pos) {
	bn.pos.EndCol = pos.EndCol
	bn.pos.EndLine = pos.EndLine
}

func (bn *BaseNode) SetTok(tok token.Token) {
	bn.Token = tok
}
