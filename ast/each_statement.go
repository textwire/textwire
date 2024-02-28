package ast

import (
	"bytes"

	"github.com/textwire/textwire/token"
)

type EachStatement struct {
	Token       token.Token // The '@for' token
	Var         *Identifier // The variable name
	Array       Expression  // The array to loop over
	Alternative *BlockStmt  // The @else block
	Block       *BlockStmt
}

func (es *EachStatement) statementNode() {
}

func (es *EachStatement) TokenLiteral() string {
	return es.Token.Literal
}

func (es *EachStatement) String() string {
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

func (es *EachStatement) Line() uint {
	return es.Token.Line
}
