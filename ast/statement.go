package ast

type Statement interface {
	Node
	statementNode()
}
