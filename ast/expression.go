package ast

type Expression interface {
	Node
	expressionNode()
}
