package ast

import (
	"github.com/textwire/textwire/token"
)

type AssignStatement struct {
	Token token.Token // The 'var' or identifier token
	Name  *Identifier
	Value Expression
}

func (as *AssignStatement) statementNode() {
}

func (as *AssignStatement) TokenLiteral() string {
	return as.Token.Literal
}

func (as *AssignStatement) String() string {
	return as.Name.String() + " = " + as.Value.String()
}

func (as *AssignStatement) Line() uint {
	return as.Token.Line
}
