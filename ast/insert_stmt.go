package ast

import (
	"bytes"
	"fmt"

	"github.com/textwire/textwire/token"
)

type InsertStmt struct {
	Token    token.Token    // The '@insert' token
	Name     *StringLiteral // The name of the insert statement
	Argument Expression     // The argument to the insert statement; nil if has block
	Block    *BlockStmt     // The block of the insert statement; nil if has argument
}

func (is *InsertStmt) statementNode() {
}

func (is *InsertStmt) TokenLiteral() string {
	return is.Token.Literal
}

func (is *InsertStmt) String() string {
	var out bytes.Buffer

	if is.Argument != nil {
		out.WriteString(fmt.Sprintf(`@insert("%s", %s)`, is.Name.String(), is.Argument.String()))
		return out.String()
	}

	out.WriteString(fmt.Sprintf(`@insert("%s")`, is.Name.String()))
	out.WriteString(is.Block.String())
	out.WriteString(`@end`)

	return out.String()
}

func (is *InsertStmt) Line() uint {
	return is.Token.Line
}
