package ast

import (
	"testing"
)

func TestFindSlotIndex(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		slots := []Statement{
			&SlotStmt{Name: &StringLiteral{Value: "country"}},
			&SlotStmt{Name: &StringLiteral{Value: "city"}},
			&SlotStmt{Name: &StringLiteral{Value: "street"}},
		}

		if idx := findSlotIndex(slots, "city"); idx != 1 {
			t.Errorf("Function should return index 1 but got %d", idx)
		}
	})

	t.Run("not found", func(t *testing.T) {
		slots := []Statement{
			&SlotStmt{Name: &StringLiteral{Value: "country"}},
			&SlotStmt{Name: &StringLiteral{Value: "city"}},
			&SlotStmt{Name: &StringLiteral{Value: "street"}},
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
		slots := []*SlotStmt{
			{Name: &StringLiteral{Value: "lastname"}},
			{Name: &StringLiteral{Value: "lastName"}},
			{Name: &StringLiteral{Value: expectDupl}},
			{Name: &StringLiteral{Value: expectDupl}},
			{Name: &StringLiteral{Value: expectDupl}},
		}

		slot, times := findDuplicateSlot(slots)
		if times != expectTimes {
			t.Fatalf("Should find %d duplicate slots, found %d", expectTimes, times)
		}

		if slot == nil {
			t.Fatalf("Function returned nil instead of slot")
		}

		if slot.Name.Value != expectDupl {
			t.Fatalf("The duplicate slot name must be %s, got %s", expectDupl, slot.Name.Value)
		}
	})

	t.Run("returns nil and 0 for no duplicates", func(t *testing.T) {
		slots := []*SlotStmt{
			{Name: &StringLiteral{Value: "lastname"}},
			{Name: &StringLiteral{Value: "lastName"}},
			{Name: &StringLiteral{Value: "last_name"}},
			{Name: &StringLiteral{Value: "last-name"}},
			{Name: &StringLiteral{Value: "LastName"}},
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
