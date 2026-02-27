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

	INC: "++",
	DEC: "--",

	NOT:    "!",
	ASSIGN: "=",
	EQ:     "==",
	NOT_EQ: "!=",

	LTHAN:    "<",
	GTHAN:    ">",
	LTHAN_EQ: "<=",
	GTHAN_EQ: ">=",

	LBRACES:  "{{",
	RBRACES:  "}}",
	LBRACE:   "{",
	RBRACE:   "}",
	LPAREN:   "(",
	RPAREN:   ")",
	LBRACKET: "[",
	RBRACKET: "]",
	DOT:      ".",
	SEMI:     ";",

	QUESTION: "?",
	COLON:    ":",
	COMMA:    ",",

	TRUE:  "true",
	FALSE: "false",
	NIL:   "nil",
	IN:    "in",

	USE:         "@use",
	RESERVE:     "@reserve",
	INSERT:      "@insert",
	FOR:         "@for",
	BREAK:       "@break",
	CONTINUE:    "@continue",
	BREAK_IF:    "@breakif",
	CONTINUE_IF: "@continueif",
	IF:          "@if",
	ELSE:        "@else",
	ELSE_IF:     "@elseif",
	END:         "@end",
	COMPONENT:   "@component",
	SLOT:        "@slot",
}

func String(t TokenType) string {
	return tokens[t]
}
