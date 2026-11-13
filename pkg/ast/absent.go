package ast

import (
	"codeberg.org/textwire/textwire/v6/pkg/position"
	"codeberg.org/textwire/textwire/v6/pkg/token"
)

type Absent struct {
	BaseNode
}

func NewAbsent(pos *position.Pos) *Absent {
	tok := token.Token{
		Type: token.ABSENT,
		Lit:  "",
		Pos:  pos,
	}

	return &Absent{NewBaseNode(tok)}
}

func (*Absent) expressionNode() {}
func (*Absent) statementNode()  {}
func (*Absent) segmentNode()    {}

func (na *Absent) String() string {
	return ""
}
