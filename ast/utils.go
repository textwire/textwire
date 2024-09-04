package ast

func findSlotStmtIndex(stmts []Statement, slotName string) int {
	for i, stmt := range stmts {
		if slot, ok := stmt.(*SlotStmt); ok && slot.Name.Value == slotName {
			return i
		}
	}

	return -1
}
