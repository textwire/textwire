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

func TestFindDuplicateProvide(t *testing.T) {
	t.Run("returns duplicate provide", func(t *testing.T) {
		expectTimes := 3
		expectDuplicate := "firstName"
		slots := []*ProvideDir{
			NewProvideDir(token.Token{}, &StrExpr{Val: "lastname"}, ""),
			NewProvideDir(token.Token{}, &StrExpr{Val: "lastName"}, ""),
			NewProvideDir(token.Token{}, &StrExpr{Val: expectDuplicate}, ""),
			NewProvideDir(token.Token{}, &StrExpr{Val: expectDuplicate}, ""),
			NewProvideDir(token.Token{}, &StrExpr{Val: expectDuplicate}, ""),
		}

		slot, times := findDuplicateProvide(slots)
		if times != expectTimes {
			t.Fatalf("Should find %d duplicate slots, found %d", expectTimes, times)
		}

		if slot == nil {
			t.Fatalf("Function returned nil instead of slot")
		}

		if slot.Name.Val != expectDuplicate {
			t.Fatalf("The duplicate slot name must be %s, got %s", expectDuplicate, slot.Name.Val)
		}
	})

	t.Run("returns nil and 0 for no duplicates", func(t *testing.T) {
		slots := []*ProvideDir{
			NewProvideDir(token.Token{}, &StrExpr{Val: "lastname"}, ""),
			NewProvideDir(token.Token{}, &StrExpr{Val: "lastName"}, ""),
			NewProvideDir(token.Token{}, &StrExpr{Val: "last_name"}, ""),
			NewProvideDir(token.Token{}, &StrExpr{Val: "last-name"}, ""),
			NewProvideDir(token.Token{}, &StrExpr{Val: "LastName"}, ""),
		}

		slot, times := findDuplicateProvide(slots)
		if times != 0 {
			t.Fatalf("Should find 0 duplicate slots, found %d", times)
		}

		if slot != nil {
			t.Fatalf("Function should return nil, got %v", slot)
		}
	})
}
