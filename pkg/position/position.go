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

// Line returns the end line position display. For multi-line tokens, showing
// the end line is much more useful then the start line. That's why we are
// using EndLine.
func (p *Pos) Line() uint {
	return p.EndLine + 1
}

// Col returns the end column position display. For long tokens, showing
// the end column is much more useful then the start column. That's why we
// are using EndCol.
func (p *Pos) Col() uint {
	return p.EndCol + 1
}
