package ast

import (
	"github.com/textwire/textwire/token"
)

type DefineStatement struct {
	Token token.Token // The 'var' or identifier token
	Name  *Identifier
	Value Expression
}

func (ds *DefineStatement) statementNode() {
}

func (ds *DefineStatement) TokenLiteral() string {
	return ds.Token.Literal
}

func (ds *DefineStatement) String() string {
	if ds.Token.Type == token.VAR {
		return "{{ var " + ds.Name.String() + " = " + ds.Value.String() + " }}"
	}

	return "{{ " + ds.Name.String() + " := " + ds.Value.String() + " }}"
}

func (df *DefineStatement) LineNum() uint {
	return df.Token.Line
}
