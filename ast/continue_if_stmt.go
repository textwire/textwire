package ast

import "github.com/textwire/textwire/v2/token"

type ContinueIfStmt struct {
	Token     token.Token // The '@continueIf' token
	Condition Expression
	Pos       token.Position
}

func (cis *ContinueIfStmt) statementNode() {
}

func (cis *ContinueIfStmt) TokenLiteral() string {
	return cis.Token.Literal
}

func (cis *ContinueIfStmt) String() string {
	return cis.Token.Literal + "(" + cis.Condition.String() + ")"
}

func (cis *ContinueIfStmt) Line() uint {
	return cis.Token.DebugLine
}

func (cis *ContinueIfStmt) Position() token.Position {
	return cis.Pos
}
