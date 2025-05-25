package ast

import (
	"fmt"

	"github.com/textwire/textwire/v2/token"
)

type FloatLiteral struct {
	BaseNode
	Value float64
}

func NewFloatLiteral(tok token.Token, val float64) *FloatLiteral {
	return &FloatLiteral{
		BaseNode: NewBaseNode(tok),
		Value:    val,
	}
}

func (fl *FloatLiteral) expressionNode() {}

func (fl *FloatLiteral) String() string {
	return fmt.Sprintf("%g", fl.Value)
}
