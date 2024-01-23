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
	DEFINE // :=

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
	LPAREN   // (
	RPAREN   // )
	LBRACKET // [
	RBRACKET // ]
	QUESTION // ?
	COLON    // :
	COMMA    // ,
	PERIOD   // .
	SEMI     // ;

	// Keywords
	IF
	ELSE
	ELSEIF
	END
	TRUE
	FALSE
	FOR
	NIL
	VAR
	LAYOUT
	RESERVE
	INSERT
	BREAK
	CONTINUE
)

var keywords = map[string]TokenType{
	"if":       IF,
	"else":     ELSE,
	"else if":  ELSEIF,
	"end":      END,
	"true":     TRUE,
	"false":    FALSE,
	"nil":      NIL,
	"var":      VAR,
	"layout":   LAYOUT,
	"reserve":  RESERVE,
	"insert":   INSERT,
	"for":      FOR,
	"break":    BREAK,
	"continue": CONTINUE,
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
