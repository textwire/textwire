package ast

import "github.com/textwire/textwire/v2/token"

type HTMLStmt struct {
	Token token.Token
	Pos   Position
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

func (hs *HTMLStmt) Position() Position {
	return hs.Pos
}
