package lexer

import (
	"strings"

	"codeberg.org/textwire/textwire/v5/pkg/token"
)

type directiveSpan struct {
	start int
	end   int
}

// skippedOffsets returns the byte offsets of lines whose only non-whitespace
// content is a single directive, including the line ending. Directive bytes are
// included, but they never reach text tokens, so only the surrounding
// whitespace is observable.
func skippedOffsets(input string) map[int]struct{} {
	skipped := map[int]struct{}{}
	dirIDs := directiveIDs(input, directiveSpans(input))

	for start := 0; start < len(input); {
		end, next := lineBounds(input, start)

		if lineHasSingleDirective(input, dirIDs, start, end) {
			markRange(skipped, start, next)
		}

		start = next
	}

	return skipped
}

// directiveIDs labels every offset with the 1-based index of the
// directive span it belongs to, or 0 when no directive covers it.
func directiveIDs(input string, spans []directiveSpan) []int {
	ids := make([]int, len(input))

	for i, span := range spans {
		for j := span.start; j < span.end && j < len(input); j++ {
			ids[j] = i + 1
		}
	}

	return ids
}

// lineBounds returns the offset where the content of the line starting at
// start ends, not counting the line ending, and the offset where the next
// line starts.
func lineBounds(input string, start int) (contentEnd, next int) {
	i := strings.IndexByte(input[start:], '\n')
	if i < 0 {
		return len(input), len(input)
	}

	newline := start + i
	contentEnd = newline

	if newline > start && input[newline-1] == '\r' {
		contentEnd--
	}

	return contentEnd, newline + 1
}

// lineHasSingleDirective reports whether every non-space character on the
// line belongs to the same directive. Embedded code, text and a second
// directive all make the line non-standalone.
func lineHasSingleDirective(input string, dirIDs []int, start, end int) bool {
	id := 0

	for i := start; i < end; i++ {
		if isInlineSpace(input[i]) {
			continue
		}

		switch {
		case dirIDs[i] == 0:
			return false
		case id == 0:
			id = dirIDs[i]
		case id != dirIDs[i]:
			return false
		}
	}

	return id != 0
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

		span, next := scanDirective(l, tok)
		spans = append(spans, span)
		tok = next
	}

	return spans
}

// scanDirective returns the span of the directive tok, including its
// parenthesized arguments, and the token after it.
func scanDirective(l *Lexer, tok token.Token) (directiveSpan, token.Token) {
	span := directiveSpan{
		start: l.pos - len(tok.Lit),
		end:   l.pos,
	}

	tok = l.Next()
	if tok.Type == token.LPAREN {
		span.end, tok = skipDirectiveArgs(l)
	}

	return span, tok
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

func markRange(skipped map[int]struct{}, from, to int) {
	for i := from; i < to; i++ {
		skipped[i] = struct{}{}
	}
}

func isInlineSpace(ch byte) bool {
	return ch == ' ' || ch == '\t'
}
