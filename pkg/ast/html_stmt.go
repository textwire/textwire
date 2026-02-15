package ast

import "github.com/textwire/textwire/v3/pkg/token"

// HTMLStmt holds literal string of HTML code. The token literal
// is its value.
type HTMLStmt struct {
	BaseNode
}

func NewHTMLStmt(tok token.Token) *HTMLStmt {
	return &HTMLStmt{
		BaseNode: NewBaseNode(tok),
	}
}

func (hs *HTMLStmt) statementNode() {}

func (hs *HTMLStmt) String() string {
	return hs.Token.Literal
}
