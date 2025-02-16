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
	':': token.COLON,
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
	// The input string to be tokenized.
	input string

	// Zero-based current character position in the input.
	pos int

	// Zero-based next character position in the input.
	nextPos int

	// Current byte character in the input.
	char byte

	// Starts from 1. Increments when a new line is found.
	// Is shown error messages. Don't confuse with lineIndex.
	debugLine uint

	// Zero-based current column index on the line.
	colIndex uint

	// Zero-based current index on the line.
	lineIndex uint

	// Determines if current character is in HTML or Textwire.
	isHTML bool

	// Determines if current character is a part of directive.
	isDirective bool

	// We increment it when we find "(" and decrement when we find ")".
	// It helps to determine if we are lexing a directive.
	countDirectiveParentheses int

	// If this is 0 and we find "}}" then it's the closing token.
	// We increment it when we find "{" and decrement when we find "}".
	// It helps to determine if we are in HTML or Textwire.
	countCurlyBraces int
}

func New(input string) *Lexer {
	l := &Lexer{
		input:       input,
		debugLine:   1,
		isHTML:      true,
		isDirective: false,
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
		return l.newToken(token.EOF, "", l.lineIndex, l.colIndex)
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

	startLine, startCol := l.lineIndex, l.colIndex
	return l.newToken(token.HTML, l.readHTML(), startLine, startCol)
}

func (l *Lexer) bracesToken(tok token.TokenType, literal string) token.Token {
	l.isHTML = tok != token.LBRACES

	startLine, startCol := l.lineIndex, l.colIndex
	l.advanceChar() // skip first brace
	l.advanceChar() // skip second brace

	return l.newToken(tok, literal, startLine, startCol)
}

func (l *Lexer) illegalToken() token.Token {
	return l.newToken(token.ILLEGAL, string(l.char), l.lineIndex, l.colIndex)
}

func (l *Lexer) directiveToken() token.Token {
	if l.char != '@' {
		return l.illegalToken()
	}

	startLine, startCol := l.lineIndex, l.colIndex
	tok, keyword := l.readDirective()

	if tok == token.ILLEGAL {
		return l.illegalToken()
	}

	hasOptionalParens := tokensWithOptionalParens[tok] && l.char == '('
	hasNoParens := tokensWithoutParens[tok]

	l.isDirective = hasOptionalParens || !hasNoParens
	l.isHTML = !l.isDirective

	return l.newToken(tok, keyword, startLine, startCol)
}

func (l *Lexer) embeddedCodeToken() token.Token {
	// check simple tokens first
	if tok, ok := simpleTokens[l.char]; ok {
		return l.newTokenAndAdvance(tok, string(l.char), l.lineIndex, l.colIndex)
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
		startLine, startCol := l.lineIndex, l.colIndex
		str := l.readString()
		return l.newTokenAndAdvance(token.STR, str, startLine, startCol)
	case '<':
		if l.peekChar() == '=' {
			startLine, startCol := l.lineIndex, l.colIndex
			l.advanceChar() // skip "="
			return l.newTokenAndAdvance(token.LTHAN_EQ, "<=", startLine, startCol)
		}

		return l.newTokenAndAdvance(token.LTHAN, "<", l.lineIndex, l.colIndex)
	case '>':
		if l.peekChar() == '=' {
			startLine, startCol := l.lineIndex, l.colIndex
			l.advanceChar() // skip "="
			return l.newTokenAndAdvance(token.GTHAN_EQ, ">=", startLine, startCol)
		}

		return l.newTokenAndAdvance(token.GTHAN, ">", l.lineIndex, l.colIndex)
	case '!':
		if l.peekChar() == '=' {
			startLine, startCol := l.lineIndex, l.colIndex
			l.advanceChar() // skip "="
			return l.newTokenAndAdvance(token.NOT_EQ, "!=", startLine, startCol)
		}

		return l.newTokenAndAdvance(token.NOT, "!", l.lineIndex, l.colIndex)
	case '-':
		if l.peekChar() == '-' {
			startLine, startCol := l.lineIndex, l.colIndex
			l.advanceChar() // skip "-"
			return l.newTokenAndAdvance(token.DEC, "--", startLine, startCol)
		}

		return l.newTokenAndAdvance(token.SUB, "-", l.lineIndex, l.colIndex)
	case '+':
		if l.peekChar() == '+' {
			startLine, startCol := l.lineIndex, l.colIndex
			l.advanceChar() // skip "+"
			return l.newTokenAndAdvance(token.INC, "++", startLine, startCol)
		}

		return l.newTokenAndAdvance(token.ADD, "+", l.lineIndex, l.colIndex)
	case '=':
		if l.peekChar() == '=' {
			startLine, startCol := l.lineIndex, l.colIndex
			l.advanceChar() // skip "="
			return l.newTokenAndAdvance(token.EQ, "==", startLine, startCol)
		}

		return l.newTokenAndAdvance(token.ASSIGN, "=", l.lineIndex, l.colIndex)
	}

	if isIdent(l.char) {
		startLine, startCol := l.lineIndex, l.colIndex
		ident := l.readIdentifier()
		return l.newToken(token.LookupIdent(ident), ident, startLine, startCol)
	}

	if isNumber(l.char) {
		startLine, startCol := l.lineIndex, l.colIndex
		num, isInt := l.readNumber()

		if isInt {
			return l.newToken(token.INT, num, startLine, startCol)
		}

		return l.newToken(token.FLOAT, num, startLine, startCol)
	}

	return l.newToken(token.ILLEGAL, string(l.char), l.lineIndex, l.colIndex)
}

func (l *Lexer) leftBraceToken() token.Token {
	l.countCurlyBraces += 1
	return l.newTokenAndAdvance(token.LBRACE, "{", l.lineIndex, l.colIndex)
}

func (l *Lexer) rightBraceToken() token.Token {
	l.countCurlyBraces -= 1
	return l.newTokenAndAdvance(token.RBRACE, "}", l.lineIndex, l.colIndex)
}

func (l *Lexer) leftParenthesesToken() token.Token {
	if l.isDirective {
		l.countDirectiveParentheses += 1
	}

	return l.newTokenAndAdvance(token.LPAREN, "(", l.lineIndex, l.colIndex)
}

func (l *Lexer) rightParenthesesToken() token.Token {
	if l.isDirective {
		l.countDirectiveParentheses -= 1
	}

	if l.isDirective && l.countDirectiveParentheses == 0 {
		l.isDirective = false
		l.isHTML = true
	}

	return l.newTokenAndAdvance(token.RPAREN, ")", l.lineIndex, l.colIndex)
}

func (l *Lexer) newToken(tokType token.TokenType, literal string, startLine, startCol uint) token.Token {
	pos := token.Position{
		StartCol:  startCol,
		EndCol:    l.colIndex,
		StartLine: startLine,
		EndLine:   l.lineIndex,
	}

	return token.Token{
		Type:      tokType,
		Literal:   literal,
		DebugLine: l.debugLine,
		Pos:       pos,
	}
}

func (l *Lexer) newTokenAndAdvance(tokType token.TokenType, literal string, startLine, startCol uint) token.Token {
	tok := l.newToken(tokType, literal, startLine, startCol)
	l.advanceChar()

	return tok
}

func (l *Lexer) readIdentifier() string {
	pos := l.pos

	for isIdent(l.char) || isNumber(l.char) {
		l.advanceChar()
	}

	return l.input[pos:l.pos]
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
		if l.pos+i > len(l.input) {
			return false
		}

		keyword := l.input[l.pos : l.pos+i]

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

	pos := l.pos

	for {
		prevChar := l.char

		l.advanceChar()

		if l.char == quote && prevChar != '\\' {
			break
		}
	}

	result = l.input[pos:l.pos]

	// remove slashes before quotes
	return strings.ReplaceAll(result, "\\"+string(quote), string(quote))
}

func (l *Lexer) readNumber() (string, bool) {
	pos := l.pos
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

	return l.input[pos:l.pos], isInt
}

func (l *Lexer) readHTML() string {
	var out bytes.Buffer

	for l.isHTML && l.char != 0 {
		if l.peekChar() == '{' && l.char != '\\' {
			break
		}

		if l.isNewLine() {
			l.advanceLine()
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
	if l.nextPos >= len(l.input) {
		l.char = 0
	} else {
		l.char = l.input[l.nextPos]
	}

	l.pos = l.nextPos
	l.nextPos += 1
}

func (l *Lexer) skipWhitespace() {
	for l.char == ' ' || l.char == '\t' || l.char == '\n' || l.char == '\r' {
		if l.isNewLine() {
			l.advanceLine()
		}

		l.advanceChar()
	}
}

func (l *Lexer) isNewLine() bool {
	return l.char == '\n'
}

func (l *Lexer) advanceLine() {
	l.debugLine += 1
	l.lineIndex += 1
	l.colIndex = 0
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
	if l.nextPos >= len(l.input) {
		return 0
	}

	return l.input[l.nextPos]
}
