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
	counts := map[string]int{}

	for _, slot := range slots {
		counts[slot.Name.Value]++
	}

	// find the first slot name that has a count greater than 1
	for name, times := range counts {
		if times > 1 {
			return name, times
		}
	}

	return "", 0
}
