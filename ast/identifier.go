package ast

import "github.com/textwire/textwire/v2/token"

type Identifier struct {
	Token token.Token
	Value string
	Pos   token.Position
}

func (i *Identifier) expressionNode() {
}

func (i *Identifier) Tok() *token.Token {
	return &i.Token
}

func (i *Identifier) String() string {
	return i.Value
}

func (i *Identifier) Line() uint {
	return i.Token.ErrorLine()
}

func (i *Identifier) Position() token.Position {
	return i.Pos
}
