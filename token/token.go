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

	// Operators
	ADD // +
	SUB // -
	MUL // *
	DIV // /
	MOD // %

	INC // ++
	DEC // --

	NOT    // !
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
	BREAK_IF
	CONTINUE_IF
	INSERT
	RESERVE
	BREAK
	CONTINUE
	COMPONENT
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
	"@continueIf": CONTINUE_IF,
	"@break":      BREAK,
	"@breakIf":    BREAK_IF,
	"@component":  COMPONENT,
}

type Token struct {
	Type    TokenType
	Literal string
	Line    uint
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
