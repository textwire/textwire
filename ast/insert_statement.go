package ast

import (
	"bytes"
	"fmt"

	"github.com/textwire/textwire/token"
)

type InsertStatement struct {
	Token    token.Token     // The 'insert' token
	Name     *StringLiteral  // The name of the insert statement
	Argument Expression      // The argument to the insert statement; nil if has block
	Block    *BlockStatement // The block of the insert statement; nil if has argument
}

func (is *InsertStatement) statementNode() {
}

func (is *InsertStatement) TokenLiteral() string {
	return is.Token.Literal
}

func (is *InsertStatement) String() string {
	var result bytes.Buffer

	if is.Argument != nil {
		result.WriteString(fmt.Sprintf(`@insert("%s", %s)`, is.Name.String(), is.Argument.String()))
		return result.String()
	}

	result.WriteString(fmt.Sprintf(`@insert("%s")`, is.Name.String()))
	result.WriteString(is.Block.String())
	result.WriteString(`@end`)

	return result.String()
}
