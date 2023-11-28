package lexer

import "github.com/go-temp/go-temp/token"

type Lexer struct {
	input        string
	position     int
	nextPosition int
	char         byte
	line         uint
	isHtml       bool
}

func New(input string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		isHtml: true,
	}

	l.advanceChar()

	return l
}

func (l *Lexer) NextToken() token.Token {
	if !l.isHtml {
		l.skipWhitespace()
	}

	if l.char == 0 {
		return l.newToken(token.EOF, "")
	}

	if l.char == '{' && l.peekChar() == '{' {
		l.isHtml = false
		l.advanceChar()
		l.advanceChar()
		return l.newToken(token.OPEN_BRACES, "{{")
	}

	if l.char == '}' && l.peekChar() == '}' {
		l.isHtml = true
		l.advanceChar()
		l.advanceChar()
		return l.newToken(token.CLOSE_BRACES, "}}")
	}

	if l.isHtml {
		tok := l.newToken(token.HTML, l.readHtml())
		l.advanceChar()
		return tok
	}

	return l.readEmbeddedCodeToken()
}

func (l *Lexer) readEmbeddedCodeToken() token.Token {
	switch l.char {
	case '+':
		return l.newToken(token.PLUS, "+")
	}

	if isIdent(l.char) {
		return l.readIdentifier()
	}

	if isNumber(l.char) {
		return l.readNumber()
	}

	return l.newToken(token.ILLEGAL, string(l.char))
}

func (l *Lexer) newToken(tokType token.TokenType, literal string) token.Token {
	return token.Token{
		Type:    tokType,
		Literal: literal,
		Line:    l.line,
	}
}

func (l *Lexer) readIdentifier() token.Token {
	position := l.position

	for isIdent(l.char) {
		l.advanceChar()
	}

	return l.newToken(token.IDENT, l.input[position:l.position])
}

func (l *Lexer) readNumber() token.Token {
	position := l.position

	for isNumber(l.char) {
		l.advanceChar()
	}

	return l.newToken(token.INT, l.input[position:l.position])
}

// todo: refactor readHtml to be more efficient
// make it similar to a readIdentifier method
func (l *Lexer) readHtml() string {
	var result string

	for l.char != 0 && l.isHtml && (l.char != '{' && l.peekChar() == '{') {
		if l.char == '\n' {
			l.line += 1
		}

		result += string(l.char)

		l.advanceChar()
	}

	return result
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

func (l *Lexer) skipWhitespace() {
	for l.char == ' ' || l.char == '\t' || l.char == '\n' || l.char == '\r' {
		if l.char == '\n' {
			l.line += 1
		}

		l.advanceChar()
	}
}

func (l *Lexer) peekChar() byte {
	if l.nextPosition >= len(l.input) {
		return 0
	}

	return l.input[l.nextPosition]
}
