package position

// Pos lines and colums are 0-based for compatibility with LSP protocol.
type Pos struct {
	StartLine uint
	StartCol  uint
	EndLine   uint
	EndCol    uint
}

func (p Pos) Contains(line uint, col uint) bool {
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

// Line returns the end line position display.
func (p *Pos) Line() uint {
	return p.EndLine + 1
}

// Col returns the end column position display.
func (p *Pos) Col() uint {
	return p.EndCol + 1
}
