package ast

import "github.com/textwire/textwire/v3/pkg/token"

// TextStmt holds literal string of text. The token literal is its value.
type TextStmt struct {
	BaseNode
}

func NewTextStmt(tok token.Token) *TextStmt {
	return &TextStmt{
		BaseNode: NewBaseNode(tok),
	}
}

func (hs *TextStmt) statementNode() {}

func (hs *TextStmt) String() string {
	return hs.Token.Lit
}
