package token

type Position struct {
	StartLine uint
	StartCol  uint
	EndLine   uint
	EndCol    uint
}

func (p Position) Contains(line uint, col uint) bool {
	// Line is out of range
	if line < p.StartLine || line > p.EndLine {
		return false
	}

	// Before start column on start line
	if line == p.StartLine && col < p.StartCol {
		return false
	}

	// After end column on end line
	if line == p.EndLine && col > p.EndCol {

		return false
	}

	return true
}
