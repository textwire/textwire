package ast

import (
	"fmt"

	"github.com/textwire/textwire/v3/pkg/token"
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
	if _, ok := rd.Fallback.(*Empty); ok {
		return fmt.Sprintf(`@reserve("%s")`, rd.Name)
	}
	return fmt.Sprintf(`@reserve("%s", %s)`, rd.Name, rd.Fallback)
}
