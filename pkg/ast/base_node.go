package ast

import (
	"github.com/textwire/textwire/v3/pkg/position"
	"github.com/textwire/textwire/v3/pkg/token"
)

// BaseNode is the main base for all the AST nodes.
type BaseNode struct {
	Token token.Token
	Pos   position.Pos
}

func NewBaseNode(tok token.Token) BaseNode {
	return BaseNode{
		Token: tok,
		Pos:   tok.Pos,
	}
}

func (bn *BaseNode) Line() uint {
	return bn.Token.Line()
}

func (bn *BaseNode) Tok() *token.Token {
	return &bn.Token
}

func (bn *BaseNode) TokPos() position.Pos {
	return bn.Pos
}

func (bn *BaseNode) SetEndPosition(pos position.Pos) {
	bn.Pos.EndCol = pos.EndCol
	bn.Pos.EndLine = pos.EndLine
}

func (bn *BaseNode) SetTok(tok token.Token) {
	bn.Token = tok
}
