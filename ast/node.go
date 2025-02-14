package ast

type Node interface {
	TokenLiteral() string
	String() string
	Line() uint
	Position() token.Position
}
