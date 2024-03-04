package lexer

func isIdent(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || ch == '_'
}

func isLetterWord(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || ch == '@'
}

func isNumber(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
