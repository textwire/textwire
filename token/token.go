package token

type TokenType int

const (
	// Special types
	ILLEGAL TokenType = iota // An illegal token
	EOF                      // The end of the file
	HTML                     // HTML code

	// Identifiers + literals
	IDENT // foo, bar
	INT   // 4, 24

	// Delimiters
	OPEN_BRACES  // {{
	CLOSE_BRACES // }}

	// Keywords
	IF

	// Operators
	PLUS     // +
	MINUS    // -
	ASTERISK // *
	SLASH    // /
	PERCENT  // %
)

var keywords = map[string]TokenType{
	"if": IF,
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
