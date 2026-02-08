package ast

import "github.com/textwire/textwire/v3/fail"

func FindProg(name string, programs []*Program) *Program {
	for i := range programs {
		if programs[i].Name == name {
			return programs[i]
		}
	}

	return nil
}

func CheckUndefinedInserts(prog *Program, inserts map[string]*InsertStmt) *fail.Error {
	for name := range inserts {
		if _, ok := prog.Reserves[name]; ok {
			continue
		}

		line := inserts[name].Line()
		path := inserts[name].AbsPath
		name := inserts[name].Name.Value

		return fail.New(line, path, "parser", fail.ErrAddMatchingReserve, name, name)
	}

	return nil
}

func findSlotIndex(stmts []Statement, slotName string) int {
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

func findDuplicateSlot(slots []*SlotStmt) (*SlotStmt, int) {
	counts := map[string]int{}
	firstSeen := map[string]*SlotStmt{}

	var maxSlot *SlotStmt
	var maxCount int

	for _, slot := range slots {
		name := slot.Name.Value
		counts[name]++

		if firstSeen[name] == nil {
			firstSeen[name] = slot
		}

		if counts[name] > 1 && counts[name] > maxCount {
			maxCount = counts[name]
			maxSlot = firstSeen[name]
		}
	}

	if maxCount == 0 {
		return nil, 0
	}

	return maxSlot, maxCount
}
