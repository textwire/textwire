package lexer

import (
	"testing"

	"github.com/textwire/textwire/token"
)

func TestHtml(t *testing.T) {
	inp := `<h2 class="container">The winter is good!</h2>`

	TokenizeString(t, inp, []token.Token{
		{Type: token.HTML, Literal: `<h2 class="container">The winter is good!</h2>`},
		{Type: token.EOF, Literal: ""},
	})
}

func TestIntegers(t *testing.T) {
	inp := "<div>{{ 0 1 2 3 4 5 6 7 8 9 234 41 }}</div>"

	TokenizeString(t, inp, []token.Token{
		{Type: token.HTML, Literal: "<div>"},
		{Type: token.LBRACES, Literal: "{{"},
		{Type: token.INT, Literal: "0"},
		{Type: token.INT, Literal: "1"},
		{Type: token.INT, Literal: "2"},
		{Type: token.INT, Literal: "3"},
		{Type: token.INT, Literal: "4"},
		{Type: token.INT, Literal: "5"},
		{Type: token.INT, Literal: "6"},
		{Type: token.INT, Literal: "7"},
		{Type: token.INT, Literal: "8"},
		{Type: token.INT, Literal: "9"},
		{Type: token.INT, Literal: "234"},
		{Type: token.INT, Literal: "41"},
		{Type: token.RBRACES, Literal: "}}"},
		{Type: token.HTML, Literal: "</div>"},
		{Type: token.EOF, Literal: ""},
	})
}

func TestIdentifiers(t *testing.T) {
	inp := "{{ testVar another_var nil }}"

	TokenizeString(t, inp, []token.Token{
		{Type: token.LBRACES, Literal: "{{"},
		{Type: token.IDENT, Literal: "testVar"},
		{Type: token.IDENT, Literal: "another_var"},
		{Type: token.NIL, Literal: "nil"},
		{Type: token.RBRACES, Literal: "}}"},
		{Type: token.EOF, Literal: ""},
	})
}

func TestOperators(t *testing.T) {
	inp := "{{ 1 + 2 - 3 * 4 / 5 % 6 }}"

	TokenizeString(t, inp, []token.Token{
		{Type: token.LBRACES, Literal: "{{"},
		{Type: token.INT, Literal: "1"},
		{Type: token.PLUS, Literal: "+"},
		{Type: token.INT, Literal: "2"},
		{Type: token.MINUS, Literal: "-"},
		{Type: token.INT, Literal: "3"},
		{Type: token.ASTERISK, Literal: "*"},
		{Type: token.INT, Literal: "4"},
		{Type: token.SLASH, Literal: "/"},
		{Type: token.INT, Literal: "5"},
		{Type: token.MODULO, Literal: "%"},
		{Type: token.INT, Literal: "6"},
		{Type: token.RBRACES, Literal: "}}"},
		{Type: token.EOF, Literal: ""},
	})
}

func TestStrings(t *testing.T) {
	inp := `{{ "Anna \"and\" Serhii" }}`

	TokenizeString(t, inp, []token.Token{
		{Type: token.LBRACES, Literal: "{{"},
		{Type: token.STR, Literal: `Anna "and" Serhii`},
		{Type: token.RBRACES, Literal: "}}"},
		{Type: token.EOF, Literal: ""},
	})
}

func TokenizeString(t *testing.T, input string, expectTokens []token.Token) {
	l := New(input)

	for i, expectTok := range expectTokens {
		tok := l.NextToken()

		if tok.Literal != expectTok.Literal {
			t.Fatalf("tests[%d] - literal wrong. expected=%s, got=%s",
				i, expectTok.Literal, tok.Literal)
		}

		if tok.Type != expectTok.Type {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%d, got=%d",
				i, expectTok.Type, tok.Type)
		}
	}
}
