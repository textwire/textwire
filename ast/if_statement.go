package ast

import (
	"bytes"

	"github.com/textwire/textwire/token"
)

type IfStatement struct {
	Token        token.Token        // The '@if' token
	Condition    Expression         // The truthy condition
	Consequence  *BlockStatement    // The 'then' block
	Alternative  *BlockStatement    // The @else block
	Alternatives []*ElseIfStatement // The @elseif blocks
}

func (is *IfStatement) statementNode() {
}

func (is *IfStatement) TokenLiteral() string {
	return is.Token.Literal
}

func (is *IfStatement) String() string {
	var out bytes.Buffer

	out.WriteString("@if(" + is.Condition.String() + ")\n")

	out.WriteString(is.Consequence.String())

	for _, e := range is.Alternatives {
		out.WriteString(e.String())
	}

	if is.Alternative != nil {
		out.WriteString("@else\n")
		out.WriteString(is.Alternative.String() + "\n")
	}

	out.WriteString("@end\n")

	return out.String()
}

func (is *IfStatement) Line() uint {
	return is.Token.Line
}
