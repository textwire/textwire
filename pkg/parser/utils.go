package parser

func isWhitespace(str string) bool {
	for _, char := range str {
		if char != ' ' && char != '\t' && char != '\n' && char != '\r' {
			return false
		}
	}

	return true
}
