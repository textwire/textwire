package ast

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/textwire/textwire/v2/token"
)

type ObjectLiteral struct {
	Token token.Token           // The '{' token
	Pairs map[string]Expression // The key-value pairs
	Pos   token.Position
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
		k := fmt.Sprintf(`"%s": %s`, key, value.String())
		pairs = append(pairs, k)
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()

}

func (os *ObjectLiteral) Line() uint {
	return os.Token.DebugLine
}

func (os *ObjectLiteral) Position() token.Position {
	return os.Pos
}
