package ast

import "github.com/textwire/textwire/token"

type HTMLStmt struct {
	Token token.Token
}

func (hs *HTMLStmt) statementNode() {
}

func (hs *HTMLStmt) TokenLiteral() string {
	return hs.Token.Literal
}

func (hs *HTMLStmt) String() string {
	return hs.Token.Literal
}

func (hs *HTMLStmt) Line() uint {
	return hs.Token.Line
}
