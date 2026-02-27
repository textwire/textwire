package lexer

import (
	"testing"

	"github.com/textwire/textwire/v3/pkg/token"
)

func BenchmarkReadDirective(b *testing.B) {
	code := "@if(a) @elseif(b) @end @breakif(false) @continueif(false) @each @for @reserve @use @insert('nice', 'cool')"

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
