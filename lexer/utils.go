package lexer

func isIdent(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || ch == '_'
}

func isDirective(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ch == '@'
}

func isNumber(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
