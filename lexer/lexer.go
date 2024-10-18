package lexer

import (
	"bytes"
	"strings"

	token "github.com/textwire/textwire/v2/token"
)

var simpleTokens = map[byte]token.TokenType{
	'*': token.MUL,
	'?': token.QUESTION,
	'/': token.DIV,
	'%': token.MOD,
	',': token.COMMA,
	'[': token.LBRACKET,
	']': token.RBRACKET,
	'.': token.DOT,
	';': token.SEMI,
}

var tokensWithoutParens = map[token.TokenType]bool{
	token.ELSE:     true,
	token.END:      true,
	token.BREAK:    true,
	token.CONTINUE: true,
	token.SLOT:     true,
}

var tokensWithOptionalParens = map[token.TokenType]bool{
	token.SLOT: true,
}

type Lexer struct {
	input                     string
	position                  int
	nextPosition              int
	char                      byte
	line                      uint
	isHTML                    bool
	isDirective               bool
	countDirectiveParentheses int
	countCurlyBraces          int
}

func New(input string) *Lexer {
	l := &Lexer{
		input:                     input,
		line:                      1,
		isHTML:                    true,
		isDirective:               false,
		countDirectiveParentheses: 0,
		countCurlyBraces:          0,
	}

	// set l.char to the first character
	l.advanceChar()

	return l
}

func (l *Lexer) NextToken() token.Token {
	if !l.isHTML {
		l.skipWhitespace()
	}

	if l.char == 0 {
		return l.newToken(token.EOF, "")
	}

	if l.char == '{' && l.peekChar() == '{' {
		tok := l.bracesToken(token.LBRACES, "{{")

		if l.char == '-' && l.peekChar() == '-' {
			l.skipComment()
			return l.NextToken()
		}

		return tok
	}

	if l.char == '}' && l.peekChar() == '}' && l.countCurlyBraces == 0 {
		return l.bracesToken(token.RBRACES, "}}")
	}

	if !l.isHTML {
		return l.embeddedCodeToken()
	}

	if l.isDirectiveStmt() {
		return l.directiveToken()
	}

	return l.newToken(token.HTML, l.readHTML())
}

func (l *Lexer) bracesToken(tok token.TokenType, literal string) token.Token {
	l.isHTML = tok != token.LBRACES
	l.advanceChar() // skip "{" or "}"
	l.advanceChar() // skip "{" or "}"
	return l.newToken(tok, literal)
}

func (l *Lexer) directiveToken() token.Token {
	if l.char != '@' {
		return l.newToken(token.ILLEGAL, string(l.char))
	}

	tok, keyword := l.readDirective()

	if tok == token.ILLEGAL {
		return l.newToken(tok, string(l.char))
	}

	hasOptionalParens := tokensWithOptionalParens[tok] && l.char == '('
	hasNoParens := tokensWithoutParens[tok]

	l.isDirective = hasOptionalParens || !hasNoParens
	l.isHTML = !l.isDirective

	return l.newToken(tok, keyword)
}

func (l *Lexer) embeddedCodeToken() token.Token {
	// check simple tokens first
	if tok, ok := simpleTokens[l.char]; ok {
		return l.newTokenAndAdvance(tok, string(l.char))
	}

	switch l.char {
	case '{':
		return l.leftBraceToken()
	case '}':
		return l.rightBraceToken()
	case '(':
		return l.leftParenthesesToken()
	case ')':
		return l.rightParenthesesToken()
	case '"', '\'':
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

func (l *Lexer) leftBraceToken() token.Token {
	l.countCurlyBraces += 1
	return l.newTokenAndAdvance(token.LBRACE, "{")
}

func (l *Lexer) rightBraceToken() token.Token {
	l.countCurlyBraces -= 1
	return l.newTokenAndAdvance(token.RBRACE, "}")
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
		l.isHTML = true
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

		if !l.isPotentiallyLong(tok) && tok != token.ILLEGAL {
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

func (l *Lexer) isPotentiallyLong(tok token.TokenType) bool {
	return (tok == token.ELSE && l.char == 'i' && l.peekChar() == 'f') ||
		(tok == token.BREAK && l.char == 'I' && l.peekChar() == 'f') ||
		(tok == token.CONTINUE && l.char == 'I' && l.peekChar() == 'f')
}

func (l *Lexer) readString() string {
	quote := l.char
	result := ""

	l.advanceChar() // skip the first quote

	if l.char == quote {
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
			if !isNumber(l.peekChar()) {
				break
			}

			isInt = false
		}

		l.advanceChar()
	}

	return l.input[position:l.position], isInt
}

func (l *Lexer) readHTML() string {
	var out bytes.Buffer

	for l.isHTML && l.char != 0 {
		if l.peekChar() == '{' && l.char != '\\' {
			break
		}

		if l.char == '\n' {
			l.line += 1
		}

		if esc := l.escapeDirective(); esc != 0 {
			out.WriteByte(esc)
		}

		if esc := l.escapeStatementStart(); esc != "" {
			out.WriteString(esc)
		}

		if l.isDirectiveStmt() {
			break
		}

		out.WriteByte(l.char)

		l.advanceChar()
	}

	if l.char != 0 && l.char != '@' && l.char != '{' {
		out.WriteByte(l.char)
		l.advanceChar()
	}

	return out.String()
}

func (l *Lexer) escapeDirective() byte {
	if l.char != '\\' || l.peekChar() != '@' {
		return 0
	}

	l.advanceChar() // skip "\"

	if l.isDirectiveStmt() {
		l.advanceChar() // skip "@"
		return '@'
	}

	return '\\'
}

func (l *Lexer) escapeStatementStart() string {
	if l.char != '\\' || l.peekChar() != '{' {
		return ""
	}

	l.advanceChar() // skip "\"

	if l.peekChar() != '{' {
		return "\\"
	}

	l.advanceChar() // skip "{"
	l.advanceChar() // skip "{"

	return "{{"
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

func (l *Lexer) skipComment() {
	for {
		if l.char != '-' || l.peekChar() != '-' {
			l.advanceChar()
			continue
		}

		l.advanceChar() // skip "-"
		l.advanceChar() // skip "-"

		if l.char == '}' || l.peekChar() == '}' {
			break
		}
	}

	l.isHTML = true

	l.advanceChar() // skip "}"
	l.advanceChar() // skip "}"
}

func (l *Lexer) peekChar() byte {
	if l.nextPosition >= len(l.input) {
		return 0
	}

	return l.input[l.nextPosition]
}
