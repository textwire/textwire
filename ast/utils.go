package ast

func findSlotStmtIndex(stmts []Statement, slotName string) int {
	for i, stmt := range stmts {
		slot, isSlot := stmt.(*SlotStmt)

		if !isSlot {
			continue
		}

		if slot.Name.Value == slotName {
			return i
		}
	}

	return -1
}
