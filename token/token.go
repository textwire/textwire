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
	PLUS     // +
	MINUS    // -
	ASTERISK // *
	SLASH    // /
	MODULO   // %
	PERIOD   // .
	BANG     // !
	ASSIGN   // =

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
	QUESTION // ?
	COLON    // :
	COMMA    // ,

	// Keywords
	IF
	ELSE
	ELSEIF
	END
	TRUE
	FALSE
	NIL
)

var keywords = map[string]TokenType{
	"if":     IF,
	"else":   ELSE,
	"elseif": ELSEIF, // else if
	"end":    END,
	"true":   TRUE,
	"false":  FALSE,
	"nil":    NIL,
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
