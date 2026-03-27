package ast

import "github.com/textwire/textwire/v4/pkg/fail"

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

		return fail.New(pos, path, fail.OriginLink, fail.ErrUnusedInsertDetected, name, name)
	}

	return nil
}

func findSlotIndex(chunks []Chunk, slotName string) int {
	for i := range chunks {
		slot, ok := chunks[i].(*SlotDir)
		if !ok {
			continue
		}

		if slot.Name.Val == slotName {
			return i
		}
	}
	return -1
}

func findDuplicatePasses(passDirs []*PassDir) (*PassDir, int) {
	counts := map[string]int{}
	firstSeen := map[string]*PassDir{}

	var maxSlot *PassDir
	var maxCount int

	for _, passDir := range passDirs {
		name := passDir.Name.Val
		counts[name]++

		if firstSeen[name] == nil {
			firstSeen[name] = passDir
		}

		if counts[name] > 1 && counts[name] > maxCount {
			maxCount = counts[name]
			maxSlot = passDir
		}
	}

	if maxCount == 0 {
		return nil, 0
	}

	return maxSlot, maxCount
}
