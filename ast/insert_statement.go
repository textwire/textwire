package ast

import (
	"bytes"
	"fmt"

	"github.com/textwire/textwire/token"
)

type InsertStatement struct {
	Token token.Token
	Name  *StringLiteral
	Block *BlockStatement
}

func (is *InsertStatement) statementNode() {
}

func (is *InsertStatement) TokenLiteral() string {
	return is.Token.Literal
}

func (is *InsertStatement) String() string {
	var result bytes.Buffer

	result.WriteString(fmt.Sprintf(`{{ insert "%s" }}`, is.Name.String()))
	result.WriteString(is.Block.String())
	result.WriteString(`{{ end }}`)

	return result.String()
}
