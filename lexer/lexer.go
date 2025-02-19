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

	// Current character position in the input.
	pos int

	// Next character position in the input.
	readPos int

	// Current byte character in the input.
	char byte

	// Current column index on the line.
	col uint

	// Start column index on the line.
	startCol uint

	// Determines if we should reset the column index to 0.
	shouldResetCol bool

	// Current index on the line.
	line uint

	// Start index on the line.
	startLine uint

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
		isHTML:      true,
		isDirective: false,
	}

	// set l.char to the first character
	l.readChar()

	return l
}

func (l *Lexer) NextToken() token.Token {
	if !l.isHTML {
		l.skipWhitespace()
	}

	if l.char == 0 {
		l.tokenBegins()
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

	if isDirective, _ := l.isDirectiveToken(); isDirective {
		return l.directiveToken()
	}

	return l.newToken(token.HTML, l.readHTML())
}

func (l *Lexer) bracesToken(tok token.TokenType, literal string) token.Token {
	l.isHTML = tok != token.LBRACES

	l.tokenBegins()
	l.readChar() // skip first brace
	l.readChar() // skip second brace

	return l.newToken(tok, literal)
}

func (l *Lexer) illegalToken() token.Token {
	l.tokenBegins()
	return l.newToken(token.ILLEGAL, string(l.char))
}

func (l *Lexer) directiveToken() token.Token {
	if l.char != '@' {
		return l.illegalToken()
	}

	tok, keyword := l.readDirective()

	if tok == token.ILLEGAL {
		return l.illegalToken()
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
		c := l.char
		l.tokenBegins()
		l.readChar() // skip the l.char
		return l.newToken(tok, string(c))
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
		return l.newToken(token.STR, l.readString())
	case '<':
		if l.peekChar() == '=' {
			l.tokenBegins()
			l.readChar() // skip "<"
			l.readChar() // skip "="
			return l.newToken(token.LTHAN_EQ, "<=")
		}

		l.tokenBegins()
		l.readChar() // skip "<"
		return l.newToken(token.LTHAN, "<")
	case '>':
		if l.peekChar() == '=' {
			l.tokenBegins()
			l.readChar() // skip ">"
			l.readChar() // skip "="
			return l.newToken(token.GTHAN_EQ, ">=")
		}

		l.tokenBegins()
		l.readChar() // skip ">"
		return l.newToken(token.GTHAN, ">")
	case '!':
		if l.peekChar() == '=' {
			l.tokenBegins()
			l.readChar() // skip "="
			l.readChar() // skip "="
			return l.newToken(token.NOT_EQ, "!=")
		}

		l.tokenBegins()
		l.readChar() // skip "!"
		return l.newToken(token.NOT, "!")
	case '-':
		if l.peekChar() == '-' {
			l.tokenBegins()
			l.readChar() // skip "-"
			l.readChar() // skip "-"
			return l.newToken(token.DEC, "--")
		}

		l.tokenBegins()
		l.readChar() // skip "-"
		return l.newToken(token.SUB, "-")
	case '+':
		if l.peekChar() == '+' {
			return l.incrementToken()
		}

		return l.addToken()
	case '=':
		if l.peekChar() == '=' {
			return l.equalToken()
		}

		return l.assignToken()
	}

	if isIdent(l.char) {
		ident := l.readIdentifier()
		return l.newToken(token.LookupIdent(ident), ident)
	}

	if isNumber(l.char) {
		return l.numberToken()
	}

	return l.illegalToken()
}

func (l *Lexer) incrementToken() token.Token {
	l.tokenBegins()
	l.readChar() // skip "+"
	l.readChar() // skip "+"

	return l.newToken(token.INC, "++")
}

func (l *Lexer) addToken() token.Token {
	l.tokenBegins()
	l.readChar() // skip "+"
	return l.newToken(token.ADD, "+")
}

func (l *Lexer) assignToken() token.Token {
	l.tokenBegins()
	l.readChar() // skip "="
	return l.newToken(token.ASSIGN, "=")
}

func (l *Lexer) equalToken() token.Token {
	l.tokenBegins()
	l.readChar() // skip "="
	l.readChar() // skip "="

	return l.newToken(token.EQ, "==")
}

func (l *Lexer) numberToken() token.Token {
	num, isInt := l.readNumber()

	if isInt {
		return l.newToken(token.INT, num)
	}

	return l.newToken(token.FLOAT, num)
}

func (l *Lexer) leftBraceToken() token.Token {
	l.countCurlyBraces += 1
	l.tokenBegins()
	l.readChar() // skip "{"

	return l.newToken(token.LBRACE, "{")
}

func (l *Lexer) rightBraceToken() token.Token {
	l.countCurlyBraces -= 1
	l.tokenBegins()
	l.readChar() // skip "}"

	return l.newToken(token.RBRACE, "}")
}

func (l *Lexer) leftParenthesesToken() token.Token {
	if l.isDirective {
		l.countDirectiveParentheses += 1
	}

	l.tokenBegins()
	l.readChar() // skip "("

	return l.newToken(token.LPAREN, "(")
}

func (l *Lexer) rightParenthesesToken() token.Token {
	if l.isDirective {
		l.countDirectiveParentheses -= 1
	}

	if l.isDirective && l.countDirectiveParentheses == 0 {
		l.isDirective = false
		l.isHTML = true
	}

	l.tokenBegins()
	l.readChar() // skip ")"

	return l.newToken(token.RPAREN, ")")
}

func (l *Lexer) newToken(tokType token.TokenType, literal string) token.Token {
	// We need to subtract 1 from the column index because
	// we already read the last character and incremented
	// the column index.
	endCol := l.col - 1

	// For EOF we don't need to decrement the column index
	if tokType == token.EOF {
		endCol = l.col
	}

	pos := token.Position{
		StartCol:  l.startCol,
		EndCol:    endCol,
		StartLine: l.startLine,
		EndLine:   l.line,
	}

	return token.Token{
		Type:    tokType,
		Literal: literal,
		Pos:     pos,
	}
}

func (l *Lexer) readIdentifier() string {
	pos := l.pos

	l.tokenBegins()

	for isIdent(l.char) || isNumber(l.char) {
		l.readChar()
	}

	return l.input[pos:l.pos]
}

func (l *Lexer) readDirective() (token.TokenType, string) {
	var keyword string
	var tok token.TokenType

	l.tokenBegins()

	for isLetterWord(l.char) {
		keyword += string(l.char)

		tok = token.LookupDirective(keyword)

		l.readChar()

		if !l.isPotentiallyLong(tok) && tok != token.ILLEGAL {
			break
		}
	}

	return tok, keyword
}

func (l *Lexer) isDirectiveToken() (isDirectory bool, escapedDirectory bool) {
	if l.char != '@' {
		return false, false
	}

	pos := l.pos

	longestDir := token.LongestDirective()

	for i := 1; i <= longestDir; i++ {
		if pos+i > len(l.input) {
			return false, false
		}

		keyword := l.input[pos : pos+i]

		tok := token.LookupDirective(keyword)

		if tok == token.ILLEGAL {
			continue
		}

		if l.prevChar() == '\\' {
			return false, true
		}

		return true, false
	}

	return false, false
}

func (l *Lexer) isPotentiallyLong(tok token.TokenType) bool {
	return (tok == token.ELSE && l.char == 'i' && l.peekChar() == 'f') ||
		(tok == token.BREAK && l.char == 'I' && l.peekChar() == 'f') ||
		(tok == token.CONTINUE && l.char == 'I' && l.peekChar() == 'f')
}

func (l *Lexer) readString() string {
	quote := l.char
	result := ""

	l.tokenBegins()
	l.readChar() // skip the first quote

	if l.char == quote {
		l.readChar() // skip the last quote
		return result
	}

	pos := l.pos

	for l.char != 0 {
		prevChar := l.char

		l.readChar()

		if l.char == quote && prevChar != '\\' {
			break
		}
	}

	result = l.input[pos:l.pos]

	l.readChar() // skip the last quote

	// remove slashes before quotes
	return strings.ReplaceAll(result, "\\"+string(quote), string(quote))
}

func (l *Lexer) readNumber() (string, bool) {
	pos := l.pos
	isInt := true
	l.tokenBegins()

	for isNumber(l.char) || l.char == '.' {
		if l.char == '.' {
			if !isNumber(l.peekChar()) {
				break
			}

			isInt = false
		}

		l.readChar()
	}

	return l.input[pos:l.pos], isInt
}

func (l *Lexer) areBracesToken() (areBraces bool, escapedBraces bool) {
	braces := l.char == '{' && l.peekChar() == '{'
	escapedBraces = l.prevChar() == '\\' && braces

	return braces && l.prevChar() != '\\', escapedBraces
}

func (l *Lexer) tokenBegins() {
	l.startCol = l.col
	l.startLine = l.line
}

func (l *Lexer) readHTML() string {
	var out bytes.Buffer
	l.tokenBegins()

	for l.isHTML && l.char != 0 {
		isDirective, escapedDir := l.isDirectiveToken()
		areBraces, escapedBraces := l.areBracesToken()

		if areBraces || isDirective {
			break
		}

		if escapedDir || escapedBraces {
			out.Truncate(out.Len() - 1)
		}

		out.WriteByte(l.char)
		l.readChar()
	}

	return out.String()
}

func (l *Lexer) prevChar() byte {
	if l.pos > 0 {
		return l.input[l.pos-1]
	}

	return 0
}

func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.char = 0
	} else {
		l.char = l.input[l.readPos]
	}

	l.pos = l.readPos
	l.readPos += 1

	// we don't need to increment this column on lexer
	// initialization because it should be 0 when it starts
	if l.pos > 0 {
		l.col += 1
	}

	if l.shouldResetCol {
		l.shouldResetCol = false
		l.col = 0
		l.line += 1
	}

	l.shouldResetCol = l.char == '\n'
}

func (l *Lexer) peekChar() byte {
	if l.readPos >= len(l.input) {
		return 0
	}

	return l.input[l.readPos]
}

func (l *Lexer) skipWhitespace() {
	for l.char == ' ' || l.char == '\t' || l.char == '\n' || l.char == '\r' {
		l.readChar()
	}
}

func (l *Lexer) skipComment() {
	for l.char != 0 {
		if l.char != '-' || l.peekChar() != '-' {
			l.readChar()
			continue
		}

		l.readChar() // skip "-"
		l.readChar() // skip "-"

		if l.char == '}' || l.peekChar() == '}' {
			break
		}
	}

	l.isHTML = true

	l.readChar() // skip "}"
	l.readChar() // skip "}"
}
