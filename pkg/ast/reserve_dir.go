package ast

import (
	"fmt"

	"github.com/textwire/textwire/v6/pkg/token"
)

type ReserveDir struct {
	BaseNode
	Name     *StrExpr
	Fallback Expression // the second argument
}

func NewReserveDir(tok token.Token) *ReserveDir {
	return &ReserveDir{
		BaseNode: NewBaseNode(tok),
	}
}

func (*ReserveDir) chunkNode() {}

func (rd *ReserveDir) String() string {
	if _, ok := rd.Fallback.(*Absent); ok {
		return fmt.Sprintf(`@reserve("%s")`, rd.Name.Val)
	}
	return fmt.Sprintf(`@reserve("%s", %s)`, rd.Name.Val, rd.Fallback)
}
