package token

type Position struct {
	StartLine uint
	StartCol  uint
	EndLine   uint
	EndCol    uint
}

func (p Position) Contains(line uint, char uint) bool {
	return (line > p.StartLine || (line == p.StartLine && char >= p.StartCol)) &&
		(line < p.EndLine || (line == p.EndLine && char <= p.EndCol))
}
