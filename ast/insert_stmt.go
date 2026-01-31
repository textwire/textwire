package ast

import (
	"bytes"
	"fmt"

	"github.com/textwire/textwire/v3/token"
)

type InsertStmt struct {
	BaseNode
	Name     *StringLiteral // The name of the insert statement
	Argument Expression     // The argument to the insert statement; nil if has block
	Block    *BlockStmt     // The block of the insert statement; nil if has argument
	FilePath string         // The file path of the insert statement
}

func NewInsertStmt(tok token.Token, filePath string) *InsertStmt {
	return &InsertStmt{
		BaseNode: NewBaseNode(tok),
		FilePath: filePath,
	}
}

func (is *InsertStmt) statementNode() {}

func (is *InsertStmt) String() string {
	var out bytes.Buffer

	if is.Argument != nil {
		fmt.Fprintf(&out, `@insert("%s", %s)`, is.Name.String(), is.Argument.String())
		return out.String()
	}

	fmt.Fprintf(&out, `@insert("%s")`, is.Name.String())

	out.WriteString(is.Block.String())
	out.WriteString(`@end`)

	return out.String()
}

func (is *InsertStmt) Stmts() []Statement {
	if is.Block == nil {
		return []Statement{}
	}

	return is.Block.Stmts()
}
