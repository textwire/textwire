package ast

import (
	"bytes"

	"github.com/textwire/textwire/v2/token"
)

type EachStmt struct {
	Token       token.Token // The '@each' token
	Var         *Identifier // The variable name
	Array       Expression  // The array to loop over
	Alternative *BlockStmt  // The @else block
	Block       *BlockStmt
	Pos         token.Position
}

func (es *EachStmt) statementNode() {
}

func (es *EachStmt) Stmts() []Statement {
	return es.Block.Statements
}

func (es *EachStmt) Tok() *token.Token {
	return &es.Token
}

func (es *EachStmt) String() string {
	var out bytes.Buffer

	out.WriteString("@each(")
	out.WriteString(es.Var.String())
	out.WriteString(" in ")
	out.WriteString(es.Array.String())
	out.WriteString(")\n")
	out.WriteString(es.Block.String() + "\n")

	if es.Alternative != nil {
		out.WriteString("@else\n")
		out.WriteString(es.Alternative.String() + "\n")
	}

	out.WriteString("@end\n")

	return out.String()
}

func (es *EachStmt) Line() uint {
	return es.Token.ErrorLine()
}

func (es *EachStmt) Position() token.Position {
	return es.Pos
}
