package ast

import (
	"fmt"

	"github.com/textwire/textwire/token"
)

type LayoutStatement struct {
	Token   token.Token    // The 'layout' token
	Name    *StringLiteral // The relative path to the layout
	Program *Program
}

func (ls *LayoutStatement) statementNode() {
}

func (ls *LayoutStatement) TokenLiteral() string {
	return ls.Token.Literal
}

func (ls *LayoutStatement) String() string {
	return fmt.Sprintf(`{{ layout %s }}`, ls.Name.String())
}
