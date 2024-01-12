package ast

import (
	"github.com/textwire/textwire/token"
)

type VarStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (vs *VarStatement) statementNode() {
}

func (vs *VarStatement) TokenLiteral() string {
	return vs.Token.Literal
}

func (vs *VarStatement) String() string {
	return "{{ var " + vs.Name.String() + " = " + vs.Value.String() + " }}"
}
