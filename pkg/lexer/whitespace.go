package lexer

import (
	"strings"

	"codeberg.org/textwire/textwire/v5/pkg/token"
)

type textwireCodeSpan struct {
	start int
	end   int
}

// skippedOffsets returns the byte offsets of whitespace that belongs to lines
// holding only textwire code. Those offsets are left out of text token
// literals.
func skippedOffsets(input string) map[int]struct{} {
	skipped := map[int]struct{}{}
	covered := codeCoverage(input, textwireCodeSpans(input))

	for start := 0; start < len(input); {
		end, next := lineBounds(input, start)

		if lineIsCodeOnly(input, covered, start, end) {
			markCodeOnlyLine(input, skipped, start, end, next)
		}

		start = next
	}

	return skipped
}

// codeCoverage marks every offset that belongs to a textwire code span.
func codeCoverage(input string, spans []textwireCodeSpan) []bool {
	covered := make([]bool, len(input))

	for _, span := range spans {
		for i := span.start; i < span.end && i < len(input); i++ {
			covered[i] = true
		}
	}

	return covered
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

// lineIsCodeOnly reports whether the line holds textwire code and nothing but
// spaces and tabs besides it.
func lineIsCodeOnly(input string, covered []bool, start, end int) bool {
	hasCode := false

	for i := start; i < end; i++ {
		if isInlineSpace(input[i]) {
			continue
		}

		if !covered[i] {
			return false
		}

		hasCode = true
	}

	return hasCode
}

// markCodeOnlyLine marks the indentation before the code and the spaces and
// line ending after it.
func markCodeOnlyLine(input string, skipped map[int]struct{}, start, end, next int) {
	line := input[start:end]
	first := start + len(line) - len(strings.TrimLeft(line, " \t"))
	last := start + len(strings.TrimRight(line, " \t"))

	markRange(skipped, start, first)
	markRange(skipped, last, next)
}

// textwireCodeSpans returns the source span of every directive, including its
// parenthesized arguments, and of every {{ }} or {!! !!} embedded code block.
func textwireCodeSpans(input string) []textwireCodeSpan {
	var spans []textwireCodeSpan

	l := newLexer(input)
	tok := l.Next()

	for tok.Type != token.EOF {
		var span textwireCodeSpan

		switch {
		case tokenIsDirective(tok):
			span, tok = directiveSpan(l, tok)
		case tok.Type == token.LBRACES || tok.Type == token.LBRACESRAW:
			span, tok = embeddedCodeSpan(l, tok)
		default:
			tok = l.Next()
			continue
		}

		spans = append(spans, span)
	}

	return spans
}

// directiveSpan returns the span of the directive tok, including its
// parenthesized arguments, and the token after it.
func directiveSpan(l *Lexer, tok token.Token) (textwireCodeSpan, token.Token) {
	span := textwireCodeSpan{
		start: l.pos - len(tok.Lit),
		end:   l.pos,
	}

	tok = l.Next()
	if tok.Type == token.LPAREN {
		span.end, tok = skipDirectiveArgs(l)
	}

	return span, tok
}

// embeddedCodeSpan returns the span of a {{ }} or {!! !!} block that opens
// with tok, and the token after it.
func embeddedCodeSpan(l *Lexer, open token.Token) (textwireCodeSpan, token.Token) {
	start := l.pos - len(open.Lit)
	end, next := skipEmbeddedBody(l, closingFor(open.Type))

	return textwireCodeSpan{start: start, end: end}, next
}

// closingFor returns the token that closes embedded code opened by open. It
// panics when open is not an embedded code opening token.
func closingFor(open token.TokenType) token.TokenType {
	switch open {
	case token.LBRACES:
		return token.RBRACES
	case token.LBRACESRAW:
		return token.RBRACESRAW
	}

	panic("unexpected embedded code opening token " + token.String(open))
}

// skipEmbeddedBody consumes tokens up to the closing brace of embedded code.
// It returns the end offset and the token after it.
func skipEmbeddedBody(l *Lexer, closing token.TokenType) (int, token.Token) {
	for {
		tok := l.Next()

		if tok.Type == closing || tok.Type == token.EOF {
			end := l.pos
			return end, l.Next()
		}
	}
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
