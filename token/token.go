package token

type TokenType int

const (
	// Special types
	ILLEGAL TokenType = iota // An illegal token
	EOF                      // The end of the file
	IDENT                    // foo, bar

	// Literals
	HTML // HTML code
	INT  // Integer
	STR  // String

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

func (t *Token) String() string {
	switch t.Type {
	case ILLEGAL:
		return "ILLEGAL"
	case EOF:
		return "EOF"
	case IDENT:
		return "IDENT"
	case HTML:
		return "HTML"
	case INT:
		return "INT"
	case STR:
		return "STR"
	case PLUS:
		return "PLUS"
	case MINUS:
		return "MINUS"
	case ASTERISK:
		return "ASTERISK"
	case SLASH:
		return "SLASH"
	case MODULO:
		return "MODULO"
	case PERIOD:
		return "PERIOD"
	case BANG:
		return "BANG"
	case ASSIGN:
		return "ASSIGN"
	case EQ:
		return "EQ"
	case NOT_EQ:
		return "NOT_EQ"
	case LTHAN:
		return "LTHAN"
	case GTHAN:
		return "GTHAN"
	case LTHAN_EQ:
		return "LTHAN_EQ"
	case GTHAN_EQ:
		return "GTHAN_EQ"
	case LBRACES:
		return "LBRACES"
	case RBRACES:
		return "RBRACES"
	case LPAREN:
		return "LPAREN"
	case RPAREN:
		return "RPAREN"
	case QUESTION:
		return "QUESTION"
	case COLON:
		return "COLON"
	case COMMA:
		return "COMMA"
	case IF:
		return "IF"
	case ELSE:
		return "ELSE"
	case ELSEIF:
		return "ELSEIF"
	case END:
		return "END"
	case TRUE:
		return "TRUE"
	case FALSE:
		return "FALSE"
	case NIL:
		return "NIL"
	default:
		return "UNKNOWN"
	}
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return IDENT
}
