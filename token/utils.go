package token

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",

	IDENT: "IDENT",
	HTML:  "HTML",
	INT:   "INT",
	FLOAT: "FLOAT",
	STR:   "STR",

	ADD: "+",
	SUB: "-",
	MUL: "*",
	DIV: "/",
	MOD: "%",

	NOT:    "!",
	ASSIGN: "=",
	DEFINE: ":=",
	EQ:     "==",
	NOT_EQ: "!=",

	LTHAN:    "<",
	GTHAN:    ">",
	LTHAN_EQ: "<=",
	GTHAN_EQ: ">=",

	LBRACES:  "{{",
	RBRACES:  "}}",
	LPAREN:   "(",
	RPAREN:   ")",
	LBRACKET: "[",
	RBRACKET: "]",
	PERIOD:   ".",

	QUESTION: "?",
	COLON:    ":",
	COMMA:    ",",

	IF:       "if",
	ELSE:     "else",
	ELSEIF:   "else if",
	END:      "end",
	TRUE:     "true",
	FALSE:    "false",
	NIL:      "nil",
	VAR:      "var",
	LAYOUT:   "layout",
	RESERVE:  "reserve",
	INSERT:   "insert",
	FOR:      "for",
	BREAK:    "break",
	CONTINUE: "continue",
}

func TokenString(t TokenType) string {
	return tokens[t]
}
