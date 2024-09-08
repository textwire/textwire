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

func findDuplicateSlot(slots []*SlotStmt) (string, int) {
	counts := make(map[string]int)

	for _, slot := range slots {
		counts[slot.Name.Value]++
	}

	for name, count := range counts {
		if count > 1 {
			return name, count
		}
	}

	return "", 0
}
