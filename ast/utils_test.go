package ast

import (
	"testing"

	"github.com/textwire/textwire/v3/token"
)

func TestFindSlotStmtIndex(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		stmts := []Statement{
			&SlotStmt{Name: &StringLiteral{Value: "country"}},
			&SlotStmt{Name: &StringLiteral{Value: "city"}},
			&SlotStmt{Name: &StringLiteral{Value: "street"}},
		}

		idx := findSlotStmtIndex(stmts, "city")

		if idx != 1 {
			t.Errorf("expect index 1 but got %d", idx)
		}
	})

	t.Run("not found", func(t *testing.T) {
		stmts := []Statement{
			&SlotStmt{Name: &StringLiteral{Value: "country"}},
			&SlotStmt{Name: &StringLiteral{Value: "city"}},
			&SlotStmt{Name: &StringLiteral{Value: "street"}},
		}

		idx := findSlotStmtIndex(stmts, "name")

		if idx != -1 {
			t.Errorf("expect index -1 but got %d", idx)
		}
	})
}

func Test_findDuplicateSlot(t *testing.T) {
	t.Run("returns duplicate slot", func(t *testing.T) {
		compName := "components/actors"
		expectTimes := 3
		expectDupl := "firstName"

		tok := token.Token{Type: token.SLOT, Literal: "@slot"}
		slots := []*SlotStmt{
			NewSlotStmt(tok, NewStringLiteral(tok, "lastname"), compName, true),
			NewSlotStmt(tok, NewStringLiteral(tok, "lastName"), compName, true),
			NewSlotStmt(tok, NewStringLiteral(tok, expectDupl), compName, true),
			NewSlotStmt(tok, NewStringLiteral(tok, expectDupl), compName, true),
			NewSlotStmt(tok, NewStringLiteral(tok, expectDupl), compName, true),
		}

		slot, times := findDuplicateSlot(slots)
		if times != expectTimes {
			t.Fatalf("Should find '%d' duplicate slots, found '%d'", expectTimes, times)
		}

		if slot == nil {
			t.Fatalf("Function returned nil instead of slot")
		}

		if slot.Name.Value != expectDupl {
			t.Fatalf("The duplicate slot name must be '%s', got '%s'", expectDupl, slot.Name.Value)
		}
	})

	t.Run("returns nil and 0 for no duplicates", func(t *testing.T) {
		compName := "components/actors"

		tok := token.Token{Type: token.SLOT, Literal: "@slot"}
		slots := []*SlotStmt{
			NewSlotStmt(tok, NewStringLiteral(tok, "lastname"), compName, true),
			NewSlotStmt(tok, NewStringLiteral(tok, "lastName"), compName, true),
			NewSlotStmt(tok, NewStringLiteral(tok, "last_name"), compName, true),
			NewSlotStmt(tok, NewStringLiteral(tok, "last-name"), compName, true),
			NewSlotStmt(tok, NewStringLiteral(tok, "LastName"), compName, true),
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
