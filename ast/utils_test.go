package ast

import "testing"

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
