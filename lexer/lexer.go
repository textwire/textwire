package lexer

import (
	"github.com/go-temp/go-temp/token"
)

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
		return l.newToken(token.HTML, l.readHtml())
	}

	return l.readEmbeddedCodeToken()
}

func (l *Lexer) readEmbeddedCodeToken() token.Token {
	switch l.char {
	case '+':
		return l.newToken(token.PLUS, "+")
	}

	if isIdent(l.char) {
		ident := l.readIdentifier()
		return l.newToken(token.LookupIdent(ident), ident)
	}

	if isNumber(l.char) {
		num := l.readNumber()
		return l.newToken(token.INT, num)
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

func (l *Lexer) readIdentifier() string {
	position := l.position

	for isIdent(l.char) {
		l.advanceChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position

	for isNumber(l.char) {
		l.advanceChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) readHtml() string {
	position := l.position

	for l.isHtml && l.char != 0 && (l.char != '{' && l.peekChar() != '{') {
		if l.char == '\n' {
			l.line += 1
		}

		l.advanceChar()
	}

	if l.char != 0 {
		l.advanceChar()
	}

	return l.input[position:l.position]
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
