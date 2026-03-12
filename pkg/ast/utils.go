package ast

import "github.com/textwire/textwire/v3/pkg/fail"

func FindProg(name string, programs []*Program) *Program {
	for i := range programs {
		if programs[i].Name == name {
			return programs[i]
		}
	}
	return nil
}

func CheckUnusedInserts(prog *Program, inserts map[string]*InsertDir) *fail.Error {
	for name := range inserts {
		if _, ok := prog.Reserves[name]; ok {
			continue
		}

		pos := inserts[name].Pos()
		path := inserts[name].AbsPath
		name := inserts[name].Name.Val

		return fail.New(pos, path, "linker", fail.ErrUnusedInsertDetected, name, name)
	}

	return nil
}

func findSlotIndex(chunks []Chunk, slotName string) int {
	for i := range chunks {
		slot, ok := chunks[i].(*SlotDir)
		if !ok {
			continue
		}

		if slot.Name().Val == slotName {
			return i
		}
	}
	return -1
}

func findDuplicateSlot(slots []SlotDirective) (SlotDirective, int) {
	counts := map[string]int{}
	firstSeen := map[string]SlotDirective{}

	var maxSlot SlotDirective
	var maxCount int

	for _, slot := range slots {
		name := slot.Name().Val
		counts[name]++

		if firstSeen[name] == nil {
			firstSeen[name] = slot
		}

		if counts[name] > 1 && counts[name] > maxCount {
			maxCount = counts[name]
			maxSlot = slot
		}
	}

	if maxCount == 0 {
		return nil, 0
	}

	return maxSlot, maxCount
}
