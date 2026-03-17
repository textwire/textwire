package ast

import (
	"github.com/textwire/textwire/v4/pkg/position"
	"github.com/textwire/textwire/v4/pkg/token"
)

type Empty struct {
	BaseNode
}

func NewEmpty(pos *position.Pos) *Empty {
	tok := token.Token{
		Type: token.EMPTY,
		Lit:  "",
		Pos:  pos,
	}

	return &Empty{NewBaseNode(tok)}
}

func (*Empty) expressionNode() {}
func (*Empty) statementNode()  {}
func (*Empty) segmentNode()    {}

func (ne *Empty) String() string {
	return ""
}
