package ast

import (
	"github.com/textwire/textwire/v2/token"
)

type AssignStmt struct {
	Token token.Token // The 'var' or identifier token
	Name  *Identifier
	Value Expression
	Pos   Position
}

func (as *AssignStmt) statementNode() {
}

func (as *AssignStmt) TokenLiteral() string {
	return as.Token.Literal
}

func (as *AssignStmt) String() string {
	return as.Name.String() + " = " + as.Value.String()
}

func (as *AssignStmt) Line() uint {
	return as.Token.StartLine
}

func (as *AssignStmt) Position() Position {
	return as.Pos
}
