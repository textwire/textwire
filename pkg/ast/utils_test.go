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

func TestFindDuplicatePasses(t *testing.T) {
	t.Run("returns duplicate pass", func(t *testing.T) {
		expectTimes := 3
		expectDuplicate := "firstName"
		slots := []*PassDir{
			NewPassDir(token.Token{}, &StrExpr{Val: "lastname"}),
			NewPassDir(token.Token{}, &StrExpr{Val: "lastName"}),
			NewPassDir(token.Token{}, &StrExpr{Val: expectDuplicate}),
			NewPassDir(token.Token{}, &StrExpr{Val: expectDuplicate}),
			NewPassDir(token.Token{}, &StrExpr{Val: expectDuplicate}),
		}

		duplicate, times := findDuplicatePasses(slots)
		if times != expectTimes {
			t.Fatalf("should find %d duplicate passes, found %d", expectTimes, times)
		}

		if duplicate == nil {
			t.Fatalf("function returned nil instead of slot")
		}

		if duplicate.Name.Val != expectDuplicate {
			t.Fatalf("duplicate pass name must be %s, got %s", expectDuplicate, duplicate.Name.Val)
		}
	})

	t.Run("returns nil and 0 for no duplicates", func(t *testing.T) {
		slots := []*PassDir{
			NewPassDir(token.Token{}, &StrExpr{Val: "lastname"}),
			NewPassDir(token.Token{}, &StrExpr{Val: "lastName"}),
			NewPassDir(token.Token{}, &StrExpr{Val: "last_name"}),
			NewPassDir(token.Token{}, &StrExpr{Val: "last-name"}),
			NewPassDir(token.Token{}, &StrExpr{Val: "LastName"}),
		}

		duplicate, times := findDuplicatePasses(slots)
		if times != 0 {
			t.Fatalf("should find 0 duplicate passes, found %d", times)
		}

		if duplicate != nil {
			t.Fatalf("function should return nil, got %v", duplicate)
		}
	})
}
