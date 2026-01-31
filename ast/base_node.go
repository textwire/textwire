package ast

import "github.com/textwire/textwire/v3/token"

// BaseNode is the base for all the ast nodes
type BaseNode struct {
	Token token.Token
	Pos   token.Position
}

func NewBaseNode(tok token.Token) BaseNode {
	return BaseNode{
		Token: tok,
		Pos:   tok.Pos,
	}
}

func (bn *BaseNode) Line() uint {
	return bn.Token.ErrorLine()
}

func (bn *BaseNode) Tok() *token.Token {
	return &bn.Token
}

func (bn *BaseNode) Position() token.Position {
	return bn.Pos
}

func (bn *BaseNode) SetEndPosition(pos token.Position) {
	bn.Pos.EndCol = pos.EndCol
	bn.Pos.EndLine = pos.EndLine
}
