package lexer

import "github.com/go-temp/go-temp/token"

type Lexer struct {
	input        string
	position     int
	nextPosition int
	char         byte
	line         uint
	file         string
	isHtml       bool
}

func (l *Lexer) NextToken() token.Token {
	if l.isHtml {
		l.skipWhitespace()
	}

	l.setLineAndFile()

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

func (l *Lexer) setLineAndFile() {
}

func (l *Lexer) readEmbeddedCodeToken() token.Token {
	//
}

func (l *Lexer) newToken(tokType token.TokenType, literal string) token.Token {
	return token.Token{
		Type:    tokType,
		Literal: literal,
		Line:    l.line,
		File:    l.file,
	}
}

func (l *Lexer) readHtml() string {
	// @todo: returns a string of HTML code
	return ""
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
		l.advanceChar()
	}
}

func (l *Lexer) peekChar() byte {
	if l.nextPosition >= len(l.input) {
		return 0
	}

	return l.input[l.nextPosition]
}
