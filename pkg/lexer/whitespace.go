package lexer

import (
	"strings"

	"codeberg.org/textwire/textwire/v5/pkg/token"
)

type directiveSpan struct {
	start int
	end   int
}

// skippedOffsets returns the byte offsets of whitespace that belongs to lines
// holding only a directive. Those offsets are left out of text token literals.
func skippedOffsets(input string) map[int]struct{} {
	skipped := map[int]struct{}{}

	for _, span := range directiveSpans(input) {
		markDirectiveLine(input, span, skipped)
	}

	return skipped
}

// directiveSpans returns the source span of every directive, including its
// parenthesized arguments.
func directiveSpans(input string) []directiveSpan {
	var spans []directiveSpan

	l := newLexer(input)
	tok := l.Next()

	for tok.Type != token.EOF {
		if !tokenIsDirective(tok) {
			tok = l.Next()
			continue
		}

		span := directiveSpan{start: l.pos - len(tok.Lit), end: l.pos}
		tok = l.Next()

		if tok.Type == token.LPAREN {
			// skipDirectiveArgs returns the token after the closing paren
			span.end, tok = skipDirectiveArgs(l)
		}

		spans = append(spans, span)
	}

	return spans
}

// skipDirectiveArgs consumes tokens up to the parenthesis that closes the
// directive arguments. It returns the end offset and the token after it.
func skipDirectiveArgs(l *Lexer) (int, token.Token) {
	depth := 1

	for {
		tok := l.Next()

		switch tok.Type {
		case token.EOF:
			return l.pos, tok
		case token.LPAREN:
			depth++
		case token.RPAREN:
			depth--
			if depth == 0 {
				// capture the offset before l.Next() moves l.pos
				end := l.pos
				return end, l.Next()
			}
		}
	}
}

// tokenIsDirective reports whether tok is a directive keyword. The type
// check keeps string literals such as "@if" from matching.
func tokenIsDirective(tok token.Token) bool {
	return tok.Type != token.ILLEGAL && token.LookupDirective(tok.Lit) == tok.Type
}

// markDirectiveLine marks the indentation before the directive and the spaces
// and line ending after it, on each line where the directive is the only
// content. A directive spanning several lines starts and ends on different
// lines, so each side is checked on its own.
func markDirectiveLine(input string, span directiveSpan, skipped map[int]struct{}) {
	lineStart := skipInlineSpacesBackward(input, span.start)
	indentOnly := lineStart == 0 || input[lineStart-1] == '\n'

	lineEnd := skipInlineSpacesForward(input, span.end)
	nextLineStart, trailingOnly := lineBreakEnd(input, lineEnd)

	multiline := strings.Contains(input[span.start:span.end], "\n")
	if !multiline && !(indentOnly && trailingOnly) {
		return
	}

	if indentOnly {
		markRange(skipped, lineStart, span.start)
	}

	if trailingOnly {
		markRange(skipped, span.end, nextLineStart)
	}
}

func skipInlineSpacesBackward(input string, pos int) int {
	for pos > 0 && isInlineSpace(input[pos-1]) {
		pos--
	}
	return pos
}

func skipInlineSpacesForward(input string, pos int) int {
	for pos < len(input) && isInlineSpace(input[pos]) {
		pos++
	}
	return pos
}

// lineBreakEnd checks what follows pos. If it is a line break or the end of
// input, it returns the offset just past it and true. Otherwise the line has
// more content, so it returns false.
func lineBreakEnd(input string, pos int) (int, bool) {
	switch {
	case pos == len(input):
		return pos, true
	case strings.HasPrefix(input[pos:], "\r\n"):
		return pos + len("\r\n"), true
	case input[pos] == '\n':
		return pos + len("\n"), true
	default:
		return pos, false
	}
}

func markRange(skipped map[int]struct{}, from, to int) {
	for i := from; i < to; i++ {
		skipped[i] = struct{}{}
	}
}

func isInlineSpace(ch byte) bool {
	return ch == ' ' || ch == '\t'
}
