package ast

import "github.com/textwire/textwire/v2/token"

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
