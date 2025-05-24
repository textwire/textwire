package ast

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/textwire/textwire/v2/token"
)

type ObjectLiteral struct {
	BaseNode
	Pairs map[string]Expression // The key-value pairs
}

func NewObjectLiteral(tok token.Token) *ObjectLiteral {
	return &ObjectLiteral{
		BaseNode: NewBaseNode(tok),
	}
}

func (ol *ObjectLiteral) expressionNode() {}

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
