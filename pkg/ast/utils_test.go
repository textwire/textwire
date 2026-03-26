package ast

import (
	"testing"

	"github.com/textwire/textwire/v4/pkg/token"
)

func TestFindSlotIndex(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		slots := []Chunk{
			NewSlotDir(token.Token{}, &StrExpr{Val: "country"}, ""),
			NewSlotDir(token.Token{}, &StrExpr{Val: "city"}, ""),
			NewSlotDir(token.Token{}, &StrExpr{Val: "street"}, ""),
		}

		if idx := findSlotIndex(slots, "city"); idx != 1 {
			t.Errorf("Function should return indaex 1 but got %d", idx)
		}
	})

	t.Run("not found", func(t *testing.T) {
		slots := []Chunk{
			NewSlotDir(token.Token{}, &StrExpr{Val: "country"}, ""),
			NewSlotDir(token.Token{}, &StrExpr{Val: "city"}, ""),
			NewSlotDir(token.Token{}, &StrExpr{Val: "street"}, ""),
		}

		if idx := findSlotIndex(slots, "name"); idx != -1 {
			t.Errorf("Function should return index -1 but got %d", idx)
		}
	})
}
