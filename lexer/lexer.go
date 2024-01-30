package lexer

import (
	"bytes"
	"strings"

	"github.com/textwire/textwire/token"
)

type Lexer struct {
	input                     string
	position                  int
	nextPosition              int
	char                      byte
	line                      uint
	isHtml                    bool
	isDirective               bool
	countDirectiveParentheses int
}

func New(input string) *Lexer {
	l := &Lexer{
		input:                     input,
		line:                      1,
		isHtml:                    true,
		isDirective:               false,
		countDirectiveParentheses: 0,
	}

	// set l.char to the first character
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
		l.advanceChar() // skip "{"
		l.advanceChar() // skip "{"
		return l.newToken(token.LBRACES, "{{")
	}

	if l.char == '}' && l.peekChar() == '}' {
		l.isHtml = true
		l.advanceChar() // skip "}"
		l.advanceChar() // skip "}"
		return l.newToken(token.RBRACES, "}}")
	}

	if !l.isHtml {
		return l.embeddedCodeToken()
	}

	if l.isDirectiveStmt() {
		return l.directiveToken()
	}

	return l.newToken(token.HTML, l.readHtml())
}

func (l *Lexer) directiveToken() token.Token {
	if l.char != '@' {
		return l.newToken(token.ILLEGAL, string(l.char))
	}

	tok, keyword := l.readDirective()

	if tok == token.ILLEGAL {
		return l.newToken(tok, string(l.char))
	}

	// ELSE and END tokens don't have parentheses
	if tok == token.ELSE || tok == token.END {
		l.isDirective = false
		l.isHtml = true
	} else {
		l.isDirective = true
		l.isHtml = false
	}

	return l.newToken(tok, keyword)
}

func (l *Lexer) embeddedCodeToken() token.Token {
	switch l.char {
	case '*':
		return l.newTokenAndAdvance(token.MUL, "*")
	case '?':
		return l.newTokenAndAdvance(token.QUESTION, "?")
	case '/':
		return l.newTokenAndAdvance(token.DIV, "/")
	case '%':
		return l.newTokenAndAdvance(token.MOD, "%")
	case ',':
		return l.newTokenAndAdvance(token.COMMA, ",")
	case '(':
		return l.leftParenthesesToken()
	case ')':
		return l.rightParenthesesToken()
	case '[':
		return l.newTokenAndAdvance(token.LBRACKET, "[")
	case ']':
		return l.newTokenAndAdvance(token.RBRACKET, "]")
	case ';':
		return l.newTokenAndAdvance(token.SEMI, ";")
	case '"', '`':
		str := l.readString()
		return l.newTokenAndAdvance(token.STR, str)
	case '<':
		if l.peekChar() == '=' {
			l.advanceChar() // skip "="
			return l.newTokenAndAdvance(token.LTHAN_EQ, "<=")
		}

		return l.newTokenAndAdvance(token.LTHAN, "<")
	case '>':
		if l.peekChar() == '=' {
			l.advanceChar() // skip "="
			return l.newTokenAndAdvance(token.GTHAN_EQ, ">=")
		}

		return l.newTokenAndAdvance(token.GTHAN, ">")
	case '!':
		if l.peekChar() == '=' {
			l.advanceChar() // skip "="
			return l.newTokenAndAdvance(token.NOT_EQ, "!=")
		}

		return l.newTokenAndAdvance(token.NOT, "!")
	case '-':
		if l.peekChar() == '-' {
			l.advanceChar() // skip "-"
			return l.newTokenAndAdvance(token.DEC, "--")
		}

		return l.newTokenAndAdvance(token.SUB, "-")
	case '+':
		if l.peekChar() == '+' {
			l.advanceChar() // skip "+"
			return l.newTokenAndAdvance(token.INC, "++")
		}

		return l.newTokenAndAdvance(token.ADD, "+")
	case '=':
		if l.peekChar() == '=' {
			l.advanceChar() // skip "="
			return l.newTokenAndAdvance(token.EQ, "==")
		}

		return l.newTokenAndAdvance(token.ASSIGN, "=")
	case ':':
		if l.peekChar() == '=' {
			l.advanceChar() // skip "="
			return l.newTokenAndAdvance(token.DEFINE, ":=")
		}

		return l.newTokenAndAdvance(token.COLON, ":")
	}

	if isIdent(l.char) {
		ident := l.readIdentifier()
		return l.newToken(token.LookupIdent(ident), ident)
	}

	if isNumber(l.char) {
		num, isInt := l.readNumber()

		if isInt {
			return l.newToken(token.INT, num)
		}

		return l.newToken(token.FLOAT, num)
	}

	return l.newToken(token.ILLEGAL, string(l.char))
}

func (l *Lexer) leftParenthesesToken() token.Token {
	if l.isDirective {
		l.countDirectiveParentheses += 1
	}

	return l.newTokenAndAdvance(token.LPAREN, "(")
}

func (l *Lexer) rightParenthesesToken() token.Token {
	if l.isDirective {
		l.countDirectiveParentheses -= 1
	}

	if l.isDirective && l.countDirectiveParentheses == 0 {
		l.isDirective = false
		l.isHtml = true
	}

	return l.newTokenAndAdvance(token.RPAREN, ")")
}

func (l *Lexer) newToken(tokType token.TokenType, literal string) token.Token {
	return token.Token{
		Type:    tokType,
		Literal: literal,
		Line:    l.line,
	}
}

func (l *Lexer) newTokenAndAdvance(tokType token.TokenType, literal string) token.Token {
	tok := token.Token{
		Type:    tokType,
		Literal: literal,
		Line:    l.line,
	}

	l.advanceChar()

	return tok
}

func (l *Lexer) readIdentifier() string {
	position := l.position

	for isIdent(l.char) || isNumber(l.char) {
		l.advanceChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) readDirective() (token.TokenType, string) {
	var keyword string
	var tok token.TokenType

	for isLetterWord(l.char) {
		keyword += string(l.char)

		tok = token.LookupDirective(keyword)

		l.advanceChar()

		if !l.isPotentialElseif(tok) && tok != token.ILLEGAL {
			break
		}
	}

	return tok, keyword
}

func (l *Lexer) isDirectiveStmt() bool {
	if l.char != '@' {
		return false
	}

	longestDir := token.LongestDirective()

	for i := 1; i <= longestDir; i++ {
		if l.position+i > len(l.input) {
			return false
		}

		keyword := l.input[l.position : l.position+i]

		tok := token.LookupDirective(keyword)

		if tok == token.ILLEGAL {
			continue
		}

		return true
	}

	return false
}

func (l *Lexer) isPotentialElseif(tok token.TokenType) bool {
	return tok == token.ELSE && l.char == 'i' && l.peekChar() == 'f'
}

func (l *Lexer) readString() string {
	quote := l.char
	result := ""

	l.advanceChar() // skip the first quote

	if l.char == quote {
		l.advanceChar()
		return result
	}

	position := l.position

	for {
		prevChar := l.char

		l.advanceChar()

		if l.char == quote && prevChar != '\\' {
			break
		}
	}

	result = l.input[position:l.position]

	// remove slashes before quotes
	return strings.ReplaceAll(result, "\\"+string(quote), string(quote))
}

func (l *Lexer) readNumber() (string, bool) {
	position := l.position
	isInt := true

	for isNumber(l.char) || l.char == '.' {
		if l.char == '.' {
			isInt = false
		}

		l.advanceChar()
	}

	return l.input[position:l.position], isInt
}

func (l *Lexer) readHtml() string {
	var out bytes.Buffer

	for l.isHtml && l.char != 0 && (l.char != '{' && l.peekChar() != '{') {
		if l.char == '\n' {
			l.line += 1
		}

		if l.char == '\\' && l.peekChar() == '@' {
			l.advanceChar() // skip "\"
			l.advanceChar() // skip "@"
			out.WriteByte('@')
		}

		if l.isDirectiveStmt() {
			break
		}

		out.WriteByte(l.char)

		l.advanceChar()
	}

	if l.char != 0 && l.char != '@' {
		out.WriteByte(l.char)
		l.advanceChar()
	}

	return out.String()
}

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
