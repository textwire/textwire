package ast

import (
	"testing"

	"github.com/textwire/textwire/v3/pkg/token"
)

func TestFindSlotIndex(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		slots := []Statement{
			NewSlotStmt(token.Token{}, &StringLiteral{Value: "country"}, "", false),
			NewSlotStmt(token.Token{}, &StringLiteral{Value: "city"}, "", false),
			NewSlotStmt(token.Token{}, &StringLiteral{Value: "street"}, "", false),
		}

		if idx := findSlotIndex(slots, "city"); idx != 1 {
			t.Errorf("Function should return index 1 but got %d", idx)
		}
	})

	t.Run("not found", func(t *testing.T) {
		slots := []Statement{
			NewSlotStmt(token.Token{}, &StringLiteral{Value: "country"}, "", false),
			NewSlotStmt(token.Token{}, &StringLiteral{Value: "city"}, "", false),
			NewSlotStmt(token.Token{}, &StringLiteral{Value: "street"}, "", false),
		}

		if idx := findSlotIndex(slots, "name"); idx != -1 {
			t.Errorf("Function should return index -1 but got %d", idx)
		}
	})
}

func TestFindDuplicateSlot(t *testing.T) {
	t.Run("returns duplicate slot", func(t *testing.T) {
		expectTimes := 3
		expectDupl := "firstName"
		slots := []SlotStatement{
			NewSlotStmt(token.Token{}, &StringLiteral{Value: "lastname"}, "", false),
			NewSlotStmt(token.Token{}, &StringLiteral{Value: "lastName"}, "", false),
			NewSlotStmt(token.Token{}, &StringLiteral{Value: expectDupl}, "", false),
			NewSlotStmt(token.Token{}, &StringLiteral{Value: expectDupl}, "", false),
			NewSlotStmt(token.Token{}, &StringLiteral{Value: expectDupl}, "", false),
		}

		slot, times := findDuplicateSlot(slots)
		if times != expectTimes {
			t.Fatalf("Should find %d duplicate slots, found %d", expectTimes, times)
		}

		if slot == nil {
			t.Fatalf("Function returned nil instead of slot")
		}

		if slot.Name().Value != expectDupl {
			t.Fatalf("The duplicate slot name must be %s, got %s", expectDupl, slot.Name().Value)
		}
	})

	t.Run("returns nil and 0 for no duplicates", func(t *testing.T) {
		slots := []SlotStatement{
			NewSlotStmt(token.Token{}, &StringLiteral{Value: "lastname"}, "", false),
			NewSlotStmt(token.Token{}, &StringLiteral{Value: "lastName"}, "", false),
			NewSlotStmt(token.Token{}, &StringLiteral{Value: "last_name"}, "", false),
			NewSlotStmt(token.Token{}, &StringLiteral{Value: "last-name"}, "", false),
			NewSlotStmt(token.Token{}, &StringLiteral{Value: "LastName"}, "", false),
		}

		slot, times := findDuplicateSlot(slots)
		if times != 0 {
			t.Fatalf("Should find 0 duplicate slots, found %d", times)
		}

		if slot != nil {
			t.Fatalf("Function should return nil, got %v", slot)
		}
	})
}
