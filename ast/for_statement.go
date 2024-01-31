package ast

import (
	"bytes"

	"github.com/textwire/textwire/token"
)

type ForStatement struct {
	Token     token.Token // The 'for' token
	Init      Statement   // initialization statement; or nil
	Condition Expression  // condition; or nil
	Post      Statement   // post statement after the loop; or nil
	Body      *BlockStatement
}

func (fs *ForStatement) statementNode() {
}

func (fs *ForStatement) TokenLiteral() string {
	return fs.Token.Literal
}

func (fs *ForStatement) String() string {
	var out bytes.Buffer

	out.WriteString("{{ for ")

	// todo: finish it here

	out.WriteString("{{ end }}\n")

	return out.String()
}

func (fs *ForStatement) Line() uint {
	return fs.Token.Line
}
