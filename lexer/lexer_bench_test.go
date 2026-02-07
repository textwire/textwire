package lexer

import (
	"testing"

	"github.com/textwire/textwire/v3/token"
)

func BenchmarkReadDirective(b *testing.B) {
	code := "@if(a) @elseif(b) @end @breakIf(false) @continueIf(false) @each @for @reserve @use @insert('nice', 'cool')"

	b.ResetTimer()

	for b.Loop() {
		lexer := New(code)

		for {
			tok := lexer.NextToken()
			if tok.Type == token.EOF {
				break
			}
		}
	}
}
