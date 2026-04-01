package token

var tokens = [...]string{
	ILLEGAL: "illegal",
	EOF:     "eof",

	IDENT: "identifier",
	TEXT:  "text",
	INT:   "integer",
	FLOAT: "float",
	STR:   "string",

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

	USE:        "@use",
	RESERVE:    "@reserve",
	INSERT:     "@insert",
	FOR:        "@for",
	BREAK:      "@break",
	CONTINUE:   "@continue",
	BREAKIF:    "@breakif",
	CONTINUEIF: "@continueif",
	IF:         "@if",
	ELSE:       "@else",
	ELSEIF:     "@elseif",
	END:        "@end",
	COMPONENT:  "@component",
	SLOT:       "@slot",
	PASS:       "@pass",
	PASSIF:     "@passif",
}

func String(t TokenType) string {
	return tokens[t]
}
