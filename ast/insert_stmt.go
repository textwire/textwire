package ast

import (
	"fmt"
	"strings"

	"github.com/textwire/textwire/v3/token"
)

type InsertStmt struct {
	BaseNode
	Name     *StringLiteral // Name of the insert statement
	Argument Expression     // Argument to the insert statement; nil if has block
	Block    *BlockStmt     // Block of the insert statement; nil if has argument
	AbsPath  string         // AbsPath of the insert statement
}

func NewInsertStmt(tok token.Token, absPath string) *InsertStmt {
	return &InsertStmt{
		BaseNode: NewBaseNode(tok),
		AbsPath:  absPath,
	}
}

func (is *InsertStmt) statementNode() {}

func (is *InsertStmt) String() string {
	var out strings.Builder
	out.Grow(30)

	if is.Argument != nil {
		fmt.Fprintf(&out, `@insert("%s", %s)`, is.Name, is.Argument)
		return out.String()
	}

	fmt.Fprintf(&out, `@insert("%s")`, is.Name)
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
