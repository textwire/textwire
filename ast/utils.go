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

func findDuplicateSlot(slots []*SlotStmt) (*SlotStmt, int) {
	occurrences := make(map[string]struct {
		slot  *SlotStmt
		count int
	})

	for _, slot := range slots {
		name := slot.Name.Value
		entry := occurrences[name]
		entry.count++

		// Store the first occurrence
		if entry.slot == nil {
			entry.slot = slot
		}

		occurrences[name] = entry

		// Return as soon as we find a duplicate
		if entry.count > 1 {
			return entry.slot, entry.count
		}
	}

	return nil, 0
}
