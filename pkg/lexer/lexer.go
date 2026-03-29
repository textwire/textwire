package lexer

import (
	"bytes"
	"strings"

	"github.com/textwire/textwire/v4/pkg/position"
	"github.com/textwire/textwire/v4/pkg/token"
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

type Lexer struct {
	// input is the input string to be tokenized.
	input string

	// pos is the current character position in the input.
	pos int

	// readPos is the next character position in the input.
	readPos int

	// char is the current byte character in the input.
	char byte

	// col is the current column index on the line.
	col uint

	// prevCol is the previous column index on the line.
	prevCol uint

	// startCol is the start column index on the line.
	startCol uint

	// shouldResetCol determines if we should reset the column index to 0.
	shouldResetCol bool

	// line is the current index on the line.
	line uint

	// prevLine is the line index of the previous character.
	// NOT THE PREVIOUS LINE, BUT THE LINE OF THE PREVIOUS CHARACTER.
	prevLine uint

	// startLine is the start index on the line.
	startLine uint

	// isText determines if current character is in text or Textwire.
	isText bool

	// isDirective determines if current character is a part of directive.
	isDirective bool

	// We increment it when we find "(" and decrement when we find ")".
	// It helps to determine if we are lexing a directive.
	countDirectiveParentheses int

	// If this is 0 and we find "}}" then it's the closing token.
	// We increment it when we find "{" and decrement when we find "}".
	// It helps to determine if we are in text or Textwire.
	countCurlyBraces int
}

func New(input string) *Lexer {
	l := &Lexer{
		input:       input,
		isText:      true,
		isDirective: false,
	}

	// set l.char to the first character
	l.readChar()

	return l
}

func (l *Lexer) Next() token.Token {
	if !l.isText {
		l.skipWhitespace()
	}

	if l.char == 0 {
		l.tokenBegins()
		return l.newToken(token.EOF, "")
	}

	if l.startsWith('{', '{', '-', '-') {
		l.skipComment()
		return l.Next()
	}

	if l.startsWith('{', '{') {
		return l.bracesToken(token.LBRACES, "{{")
	}

	if l.startsWith('}', '}') && l.countCurlyBraces == 0 {
		return l.bracesToken(token.RBRACES, "}}")
	}

	if l.isDirectiveToken() {
		return l.directiveToken()
	}

	if !l.isText {
		return l.embeddedCodeToken()
	}

	return l.newToken(token.TEXT, l.readText())
}

func (l *Lexer) startsWith(chars ...byte) bool {
	if len(chars) == 0 {
		return false
	}

	if chars[0] != l.char {
		return false
	}

	for i := 1; i < len(chars); i++ {
		if l.peek(i-1) != chars[i] {
			return false
		}
	}

	return true
}

func (l *Lexer) bracesToken(tok token.TokenType, literal string) token.Token {
	l.isText = tok != token.LBRACES

	l.tokenBegins()
	l.readChars(2) // skip braces

	return l.newToken(tok, literal)
}

func (l *Lexer) illegalToken() token.Token {
	l.tokenBegins()
	tok := l.newToken(token.ILLEGAL, string(l.char))
	l.readChar()
	return tok
}

func (l *Lexer) directiveToken() token.Token {
	if l.char != '@' {
		return l.illegalToken()
	}

	tok, keyword := l.readDirective()

	if tok == token.ILLEGAL {
		return l.illegalToken()
	}

	l.isDirective = l.char == '(' || l.nextNonSpaceIs('(')
	l.isText = !l.isDirective

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
		return l.operatorToken('=', token.LTHAN_EQ, token.LTHAN, "<=", "<")
	case '>':
		return l.operatorToken('=', token.GTHAN_EQ, token.GTHAN, ">=", ">")
	case '!':
		return l.operatorToken('=', token.NOT_EQ, token.NOT, "!=", "!")
	case '-':
		return l.operatorToken('-', token.DEC, token.SUB, "--", "-")
	case '&':
		if l.peek(0) == '&' {
			return l.twoCharToken(token.AND, "&&")
		}
		fallthrough
	case '|':
		if l.peek(0) == '|' {
			return l.twoCharToken(token.OR, "||")
		}
		fallthrough
	case '+':
		return l.operatorToken('+', token.INC, token.ADD, "++", "+")
	case '=':
		return l.operatorToken('=', token.EQ, token.ASSIGN, "==", "=")
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
		l.isText = true
	}

	l.tokenBegins()
	l.readChar() // skip ")"

	return l.newToken(token.RPAREN, ")")
}

func (l *Lexer) twoCharToken(tokType token.TokenType, literal string) token.Token {
	l.tokenBegins()
	l.readChars(2)
	return l.newToken(tokType, literal)
}

// operatorToken handles the common pattern of checking for a two-char operator
// or falling back to single-char.
func (l *Lexer) operatorToken(
	second byte,
	twoCharTok, singleTok token.TokenType,
	twoCharLit, singleLit string,
) token.Token {
	if l.peek(0) == second {
		return l.twoCharToken(twoCharTok, twoCharLit)
	}

	l.tokenBegins()
	l.readChar()

	return l.newToken(singleTok, singleLit)
}

func (l *Lexer) newToken(tokType token.TokenType, literal string) token.Token {
	endCol := l.col
	endLine := l.line

	// We need to set the end column and line to the values of the previous
	// character because we already read the last character and incremented
	// the column index.
	// For EOF and ILLEGAL we don't need to decrement the column index.
	if tokType != token.EOF && tokType != token.ILLEGAL {
		endCol = l.prevCol
		endLine = l.prevLine
	}

	pos := &position.Pos{
		StartCol:  l.startCol,
		EndCol:    endCol,
		StartLine: l.startLine,
		EndLine:   endLine,
	}

	return token.Token{
		Type: tokType,
		Lit:  literal,
		Pos:  pos,
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
	var keyword strings.Builder
	tok := token.ILLEGAL

	l.tokenBegins()

	for isLetterWord(l.char) && (l.hasIfVariant(tok) || tok == token.ILLEGAL) {
		keyword.WriteByte(l.char)
		tok = token.LookupDirective(keyword.String())
		l.readChar()
	}

	return tok, keyword.String()
}

func (l *Lexer) hasDirectivePrefix() bool {
	if l.char != '@' {
		return false
	}

	pos := l.pos
	longestDir := token.LongestDirective()

	for i := 1; i <= longestDir; i++ {
		if pos+i > len(l.input) {
			return false
		}

		if token.LookupDirective(l.input[pos:pos+i]) != token.ILLEGAL {
			return true
		}
	}

	return false
}

func (l *Lexer) isDirectiveToken() bool {
	return l.prevChar() != '\\' && l.hasDirectivePrefix()
}

func (l *Lexer) isEscapedDirective() bool {
	return l.prevChar() == '\\' && l.hasDirectivePrefix()
}

func (l *Lexer) hasIfVariant(tok token.TokenType) bool {
	// Tokens that can be extended with "if" variants, like @breakif
	longTokens := map[token.TokenType]struct{}{
		token.ELSE:     {},
		token.BREAK:    {},
		token.CONTINUE: {},
		token.PASS:     {},
	}

	if _, ok := longTokens[tok]; !ok {
		return false
	}

	return l.startsWith('i', 'f')
}

func (l *Lexer) readString() string {
	quote := l.char
	strLiteral := ""

	l.tokenBegins()
	l.readChar() // skip the first quote

	if l.char == quote {
		l.readChar() // skip the last quote
		return strLiteral
	}

	pos := l.pos

	for l.char != 0 {
		prevChar := l.char

		l.readChar()

		if l.char == quote && prevChar != '\\' {
			break
		}
	}

	strLiteral = l.input[pos:l.pos]

	l.readChar() // skip the last quote

	return handleEscapeSequences(strLiteral, quote)
}

func handleEscapeSequences(s string, quote byte) string {
	// remove slashes before quotes
	s = strings.ReplaceAll(s, "\\"+string(quote), string(quote))

	// handle other escape sequences
	s = strings.ReplaceAll(s, "\\n", "\n")
	s = strings.ReplaceAll(s, "\\t", "\t")
	s = strings.ReplaceAll(s, "\\r", "\r")

	return strings.ReplaceAll(s, "\\\\", "\\")
}

func (l *Lexer) readNumber() (string, bool) {
	pos := l.pos
	isInt := true
	l.tokenBegins()

	for isNumber(l.char) {
		l.readChar()
	}

	if l.char == '.' && isNumber(l.peek(0)) {
		isInt = false
		l.readChar()

		for isNumber(l.char) {
			l.readChar()
		}
	}

	return l.input[pos:l.pos], isInt
}

func (l *Lexer) areBracesToken() (areBraces bool, escapedBraces bool) {
	braces := l.startsWith('{', '{')
	escapedBraces = l.prevChar() == '\\' && braces

	return braces && l.prevChar() != '\\', escapedBraces
}

func (l *Lexer) tokenBegins() {
	l.startCol = l.col
	l.startLine = l.line
}

func (l *Lexer) readText() string {
	var out bytes.Buffer
	l.tokenBegins()

	for l.isText && l.char != 0 {
		areBraces, escapedBraces := l.areBracesToken()

		if areBraces || l.isDirectiveToken() {
			break
		}

		if escapedBraces || l.isEscapedDirective() {
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
	l.prevLine = l.line

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
		l.prevCol = l.col
		l.col += 1
	}

	if l.shouldResetCol {
		l.shouldResetCol = false
		l.col = 0
		l.line += 1
	}

	l.shouldResetCol = l.char == '\n'
}

func (l *Lexer) readChars(chars int) {
	for range chars {
		l.readChar()
	}
}

func (l *Lexer) peek(ahead int) byte {
	if l.readPos >= len(l.input) {
		return 0
	}

	return l.input[l.readPos+ahead]
}

func (l *Lexer) isWhitespace(char byte) bool {
	return char == ' ' || char == '\t' || char == '\n' || char == '\r'
}

// nextNonSpaceIs checks if the next non-space/tab character matches
// the given character without consuming any input.
func (l *Lexer) nextNonSpaceIs(char byte) bool {
	pos := l.readPos
	for pos < len(l.input) {
		c := l.input[pos]
		if c != ' ' && c != '\t' {
			return c == char
		}
		pos++
	}
	return false
}

func (l *Lexer) skipWhitespace() {
	for l.isWhitespace(l.char) {
		l.readChar()
	}
}

func (l *Lexer) skipComment() {
	depth := 1

	l.readChars(4) // skip "{{--"

	for l.char != 0 && depth > 0 {
		if l.startsWith('{', '{', '-', '-') {
			l.readChars(4) // skip "{{--"
			depth++
			continue
		}

		if l.startsWith('-', '-', '}', '}') {
			l.readChars(4) // skip "--}}"
			depth--
			continue
		}

		l.readChar()
	}

	l.isText = true
}
