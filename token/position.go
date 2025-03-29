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

	// On start line, check if cursor is after start column
	if line == p.StartLine && col < p.StartCol {
		return false
	}

	// On end line, check if cursor is before end column
	if line == p.EndLine && col > p.EndCol {
		return false
	}

	return true
}
