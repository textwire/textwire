package ast

func findSlotStmtIndex(stmts []Statement, slotName *StringLiteral) int {
	for i, stmt := range stmts {
		slot, isSlot := stmt.(*SlotStmt)

		if !isSlot {
			continue
		}

		if slot.Name == nil && slotName == nil {
			return i
		}

		if slot.Name.Value == slotName.Value {
			return i
		}
	}

	return -1
}
