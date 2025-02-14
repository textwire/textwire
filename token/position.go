package token

type Position struct {
	StartLine uint
	StartChar uint
	EndLine   uint
	EndChar   uint
}

func (p Position) Contains(line uint, char uint) bool {
	return (line > p.StartLine || (line == p.StartLine && char >= p.StartChar)) &&
		(line < p.EndLine || (line == p.EndLine && char <= p.EndChar))
}
