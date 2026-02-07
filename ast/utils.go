package ast

func FindProg(name string, programs []*Program) *Program {
	for i := range programs {
		if programs[i].Name == name {
			return programs[i]
		}
	}

	return nil
}

func findSlotStmtIndex(stmts []Statement, slotName string) int {
	for i, stmt := range stmts {
		slot, ok := stmt.(*SlotStmt)
		if !ok {
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
