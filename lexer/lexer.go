package lexer

type Lexer struct {
	input        string
	position     int
	nextPosition int
	char         byte
}

func (l *Lexer) NextToken() {
}

// advanceChar advances the lexer's position in the input string
func (l *Lexer) advanceChar() {
	if l.nextPosition >= len(l.input) {
		l.char = 0
	} else {
		l.char = l.input[l.nextPosition]
	}

	l.position = l.nextPosition
	l.nextPosition += 1
}
