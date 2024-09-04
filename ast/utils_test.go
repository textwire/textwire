package ast

import "testing"

func TestFindSlotStmtIndex(t *testing.T) {
	t.Run("found", func(tt *testing.T) {
		stmts := []Statement{
			&SlotStmt{Name: &StringLiteral{Value: "country"}},
			&SlotStmt{Name: &StringLiteral{Value: "city"}},
			&SlotStmt{Name: &StringLiteral{Value: "street"}},
		}

		idx := findSlotStmtIndex(stmts, &StringLiteral{Value: "city"})

		if idx != 1 {
			tt.Errorf("expected index 1 but got %d", idx)
		}
	})

	t.Run("not found", func(tt *testing.T) {
		stmts := []Statement{
			&SlotStmt{Name: &StringLiteral{Value: "country"}},
			&SlotStmt{Name: &StringLiteral{Value: "city"}},
			&SlotStmt{Name: &StringLiteral{Value: "street"}},
		}

		idx := findSlotStmtIndex(stmts, &StringLiteral{Value: "name"})

		if idx != -1 {
			tt.Errorf("expected index -1 but got %d", idx)
		}
	})
}
