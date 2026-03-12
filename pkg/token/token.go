package token

import "github.com/textwire/textwire/v3/pkg/position"

type TokenType int

const (
	// Special types
	ILLEGAL TokenType = iota // An illegal token
	EOF                      // The end of the file
	IDENT                    // foo, bar

	// Literals
	TEXT  // Text
	INT   // Integer
	FLOAT // Float
	STR   // String

	// Logical Operators
	AND // &&
	OR  // ||
	NOT // !

	// Operators
	ADD // +
	SUB // -
	MUL // *
	DIV // /
	MOD // %

	INC // ++
	DEC // --

	ASSIGN // =

	// Comparison operators
	EQ       // ==
	NOT_EQ   // !=
	LTHAN    // <
	GTHAN    // >
	LTHAN_EQ // <=
	GTHAN_EQ // >=

	// Delimiters
	LBRACES  // {{
	RBRACES  // }}
	LBRACE   // {
	RBRACE   // }
	LPAREN   // (
	RPAREN   // )
	LBRACKET // [
	RBRACKET // ]
	QUESTION // ?
	COLON    // :
	COMMA    // ,
	DOT      // .
	SEMI     // ;

	// Keywords
	TRUE
	FALSE
	NIL
	IN

	// Directives
	IF
	ELSE
	ELSEIF
	END
	FOR
	USE
	EACH
	BREAKIF
	CONTINUEIF
	INSERT
	RESERVE
	BREAK
	CONTINUE
	COMPONENT
	SLOT
	SLOTIF
	DUMP
)

var keywords = map[string]TokenType{
	// Code keywords
	"true":  TRUE,
	"false": FALSE,
	"nil":   NIL,
	"in":    IN,
}

var directives = map[string]TokenType{
	"@if":         IF,
	"@else":       ELSE,
	"@elseif":     ELSEIF,
	"@end":        END,
	"@use":        USE,
	"@reserve":    RESERVE,
	"@insert":     INSERT,
	"@for":        FOR,
	"@each":       EACH,
	"@continue":   CONTINUE,
	"@break":      BREAK,
	"@component":  COMPONENT,
	"@slotif":     SLOTIF,
	"@slot":       SLOT,
	"@dump":       DUMP,
	"@continueif": CONTINUEIF,
	"@breakif":    BREAKIF,
}

func GetDirectives() map[string]TokenType {
	return directives
}

type Token struct {
	Type TokenType
	Lit  string
	Pos  *position.Pos
}

// Line returns the end line position of the token for error display.
func (t *Token) Line() uint {
	return t.Pos.Line()
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return IDENT
}

func LongestDirective() int {
	var longest int

	for dir := range directives {
		if len(dir) > longest {
			longest = len(dir)
		}
	}

	return longest
}

func LookupDirective(dir string) TokenType {
	if tok, ok := directives[dir]; ok {
		return tok
	}

	return ILLEGAL
}
