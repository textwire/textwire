package ast

import (
	"bytes"
	"strings"

	"github.com/textwire/textwire/token"
)

type ObjectLiteral struct {
	Token token.Token               // The '{' token
	Pairs map[Expression]Expression // The key-value pairs
}

func (ol *ObjectLiteral) expressionNode() {
}

func (ol *ObjectLiteral) TokenLiteral() string {
	return ol.Token.Literal
}

func (ol *ObjectLiteral) String() string {
	var out bytes.Buffer

	pairs := []string{}
	for key, value := range ol.Pairs {
		pairs = append(pairs, key.String()+": "+value.String())
	}

	out.WriteString("{ ")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString(" }")

	return out.String()

}

func (os *ObjectLiteral) Line() uint {
	return os.Token.Line
}
