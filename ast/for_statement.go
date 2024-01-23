package ast

import (
	"bytes"

	"github.com/textwire/textwire/token"
)

type ForStatement struct {
	Token     token.Token
	Init      Statement  // initialization statement; or nil
	Condition Expression // condition or nil
	Post      Statement  // post statement after the loop; or nil
	Body      *BlockStatement
}

func (fs *ForStatement) statementNode() {
}

func (fs *ForStatement) TokenLiteral() string {
	return fs.Token.Literal
}

func (fs *ForStatement) String() string {
	var result bytes.Buffer

	result.WriteString("{{ for ")

	// todo: finish it here

	result.WriteString("{{ end }}\n")

	return result.String()
}
