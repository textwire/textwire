package ast

import "github.com/textwire/textwire/v2/token"

type ContinueIfStmt struct {
	Token     token.Token // The '@continueIf' token
	Condition Expression
}

func (cis *ContinueIfStmt) statementNode() {
}

func (cis *ContinueIfStmt) TokenLiteral() string {
	return cis.Token.Literal
}

func (cis *ContinueIfStmt) String() string {
	return cis.Token.Literal
}

func (cis *ContinueIfStmt) Line() uint {
	return cis.Token.Line
}
