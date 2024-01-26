package ast

import (
	"fmt"

	"github.com/textwire/textwire/token"
)

type UseStatement struct {
	Token   token.Token    // The 'use' token
	Name    *StringLiteral // The relative path to the layout
	Program *Program
}

func (us *UseStatement) statementNode() {
}

func (us *UseStatement) TokenLiteral() string {
	return us.Token.Literal
}

func (us *UseStatement) String() string {
	return fmt.Sprintf(`{{ use %s }}`, us.Name.String())
}
