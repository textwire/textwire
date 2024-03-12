package ast

import (
	"bytes"
	"strings"

	"github.com/textwire/textwire/token"
)

type ComponentStmt struct {
	Token     token.Token // The '@component' token
	Name      *StringLiteral
	Arguments []Expression
	Block     *Program
}

func (cs *ComponentStmt) statementNode() {
}

func (cs *ComponentStmt) TokenLiteral() string {
	return cs.Token.Literal
}

func (cs *ComponentStmt) String() string {
	var out bytes.Buffer
	var args []string

	for _, arg := range cs.Arguments {
		args = append(args, arg.String())
	}

	out.WriteString("@component(")
	out.WriteString(cs.Name.String())
	out.WriteString(", ")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

func (cs *ComponentStmt) Line() uint {
	return cs.Token.Line
}
