package ast

import (
	"bytes"

	"github.com/textwire/textwire/token"
)

type IfStatement struct {
	Token        token.Token // The 'if' token
	Condition    Expression
	Consequence  *BlockStatement
	Alternative  *BlockStatement
	Alternatives []*ElseIfStatement
}

func (is *IfStatement) statementNode() {
}

func (is *IfStatement) TokenLiteral() string {
	return is.Token.Literal
}

func (is *IfStatement) String() string {
	var result bytes.Buffer

	result.WriteString("{{ if " + is.Condition.String() + " }}\n")

	result.WriteString(is.Consequence.String())

	for _, e := range is.Alternatives {
		result.WriteString(e.String())
	}

	if is.Alternative != nil {
		result.WriteString("{{ else }}\n")
		result.WriteString(is.Alternative.String() + "\n")
	}

	result.WriteString("{{ end }}\n")

	return result.String()
}
