package ast

import (
	"fmt"
	"strings"

	"github.com/textwire/textwire/v3/pkg/token"
)

type IfStmt struct {
	BaseNode
	Condition   Expression
	IfBlock     *BlockStmt // @if()<IfBlock>@end
	ElseBlock   *BlockStmt // @else<ElseBlock>@end
	ElseifStmts []Statement
}

func NewIfStmt(tok token.Token) *IfStmt {
	return &IfStmt{
		BaseNode: NewBaseNode(tok),
	}
}

func (is *IfStmt) statementNode() {}

func (is *IfStmt) String() string {
	var out strings.Builder
	out.Grow(20 + len(is.ElseifStmts)*2)

	fmt.Fprintf(&out, "@if(%s)\n%s", is.Condition, is.IfBlock)

	for _, e := range is.ElseifStmts {
		out.WriteString(e.String())
	}

	if is.ElseBlock != nil {
		out.WriteString("@else\n")
		out.WriteString(is.ElseBlock.String() + "\n")
	}

	out.WriteString("@end\n")

	return out.String()
}

func (is *IfStmt) Stmts() []Statement {
	stmts := make([]Statement, 0)
	if is.IfBlock != nil {
		stmts = append(stmts, is.IfBlock.Stmts()...)
	}

	if is.ElseBlock != nil {
		stmts = append(stmts, is.ElseBlock.Stmts()...)
	}

	for _, e := range is.ElseifStmts {
		if withStmts, ok := e.(NodeWithStatements); ok {
			stmts = append(stmts, withStmts.Stmts()...)
		}
	}

	return stmts
}
