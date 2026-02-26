package token

type TokenType int

const (
	// Special types
	ILLEGAL TokenType = iota // An illegal token
	EOF                      // The end of the file
	IDENT                    // foo, bar

	// Literals
	HTML  // HTML code
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
	ELSE_IF
	END
	FOR
	USE
	EACH
	BREAK_IF
	CONTINUE_IF
	INSERT
	RESERVE
	BREAK
	CONTINUE
	COMPONENT
	SLOT
	SLOT_IF
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
	"@elseif":     ELSE_IF,
	"@end":        END,
	"@use":        USE,
	"@reserve":    RESERVE,
	"@insert":     INSERT,
	"@for":        FOR,
	"@each":       EACH,
	"@continue":   CONTINUE,
	"@continueIf": CONTINUE_IF,
	"@break":      BREAK,
	"@breakIf":    BREAK_IF,
	"@component":  COMPONENT,
	"@slotIf":     SLOT_IF,
	"@slot":       SLOT,
	"@dump":       DUMP,
}

func GetDirectives() map[string]TokenType {
	return directives
}

type Token struct {
	Type    TokenType
	Literal string
	Pos     Position
}

// ErrorLine returns the start line position of the token.
// It is used to display the error message and starts from 1.
func (t *Token) ErrorLine() uint {
	// add 1 because StartLine starts with 0
	return t.Pos.EndLine + 1
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
