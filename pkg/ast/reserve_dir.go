package ast

import (
	"fmt"

	"github.com/textwire/textwire/v3/pkg/token"
)

type ReserveDir struct {
	BaseNode
	Name *StrExpr
	// Fallback is the second argument; nil if not present
	Fallback Expression
}

func NewReserveDir(tok token.Token) *ReserveDir {
	return &ReserveDir{
		BaseNode: NewBaseNode(tok),
	}
}

func (_ *ReserveDir) chunkNode() {}

func (rd *ReserveDir) String() string {
	if rd.Fallback == nil {
		return fmt.Sprintf(`@reserve("%s")`, rd.Name)
	}
	return fmt.Sprintf(`@reserve("%s", %s)`, rd.Name, rd.Fallback)
}
