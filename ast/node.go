package ast

import "github.com/textwire/textwire/v2/token"

type Node interface {
	TokenLiteral() string
	String() string
	Line() uint
	Position() token.Position
}
