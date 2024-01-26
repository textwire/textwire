package lexer

import (
	"fmt"
	"testing"

	"github.com/textwire/textwire/token"
)

func TokenizeString(t *testing.T, input string, expectTokens []token.Token) {
	l := New(input)

	for i, expectTok := range expectTokens {
		tok := l.NextToken()

		if tok.Literal != expectTok.Literal {
			t.Fatalf("tests[%d] - literal wrong. expected=%s, got=%s",
				i, expectTok.Literal, tok.Literal)
		}

		if tok.Type != expectTok.Type {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%s, got=%s",
				i, expectTok.Literal, tok.Literal)
		}
	}
}

func TestHtml(t *testing.T) {
	inp := `<h2 class="container">The winter is good!</h2>`

	TokenizeString(t, inp, []token.Token{
		{Type: token.HTML, Literal: `<h2 class="container">The winter is good!</h2>`},
		{Type: token.EOF, Literal: ""},
	})
}

func TestIntegers(t *testing.T) {
	inp := "<div>{{ 0 1 2 3 4 5 6 7 8 9 234 -41 }}</div>"

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
		{Type: token.SUB, Literal: "-"},
		{Type: token.INT, Literal: "41"},
		{Type: token.RBRACES, Literal: "}}"},
		{Type: token.HTML, Literal: "</div>"},
		{Type: token.EOF, Literal: ""},
	})
}

func TestFloats(t *testing.T) {
	inp := "{{ 0.12 1.1111 9. }}"

	TokenizeString(t, inp, []token.Token{
		{Type: token.LBRACES, Literal: "{{"},
		{Type: token.FLOAT, Literal: "0.12"},
		{Type: token.FLOAT, Literal: "1.1111"},
		{Type: token.FLOAT, Literal: "9."},
		{Type: token.RBRACES, Literal: "}}"},
		{Type: token.EOF, Literal: ""},
	})
}

func TestIdentifiers(t *testing.T) {
	inp := "{{ testVar another_var12 nil false !true use reserve insert }}"

	TokenizeString(t, inp, []token.Token{
		{Type: token.LBRACES, Literal: "{{"},
		{Type: token.IDENT, Literal: "testVar"},
		{Type: token.IDENT, Literal: "another_var12"},
		{Type: token.NIL, Literal: "nil"},
		{Type: token.FALSE, Literal: "false"},
		{Type: token.NOT, Literal: "!"},
		{Type: token.TRUE, Literal: "true"},
		{Type: token.USE, Literal: "use"},
		{Type: token.RESERVE, Literal: "reserve"},
		{Type: token.INSERT, Literal: "insert"},
		{Type: token.RBRACES, Literal: "}}"},
		{Type: token.EOF, Literal: ""},
	})
}

func TestIfStatement(t *testing.T) {
	inp := "{{ if true }}1{{ else if false }}2{{ else }}3{{ end }}"

	TokenizeString(t, inp, []token.Token{
		{Type: token.LBRACES, Literal: "{{"},
		{Type: token.IF, Literal: "if"},
		{Type: token.TRUE, Literal: "true"},
		{Type: token.RBRACES, Literal: "}}"},
		{Type: token.HTML, Literal: "1"},
		{Type: token.LBRACES, Literal: "{{"},
		{Type: token.ELSEIF, Literal: "else if"},
		{Type: token.FALSE, Literal: "false"},
		{Type: token.RBRACES, Literal: "}}"},
		{Type: token.HTML, Literal: "2"},
		{Type: token.LBRACES, Literal: "{{"},
		{Type: token.ELSE, Literal: "else"},
		{Type: token.RBRACES, Literal: "}}"},
		{Type: token.HTML, Literal: "3"},
		{Type: token.LBRACES, Literal: "{{"},
		{Type: token.END, Literal: "end"},
		{Type: token.RBRACES, Literal: "}}"},
		{Type: token.EOF, Literal: ""},
	})
}

func TestOperators(t *testing.T) {
	inp := "{{ 1 + 2 - 3 * 4 / 5 % (6) 3++ + 2-- }}"

	TokenizeString(t, inp, []token.Token{
		{Type: token.LBRACES, Literal: "{{"},
		{Type: token.INT, Literal: "1"},
		{Type: token.ADD, Literal: "+"},
		{Type: token.INT, Literal: "2"},
		{Type: token.SUB, Literal: "-"},
		{Type: token.INT, Literal: "3"},
		{Type: token.MUL, Literal: "*"},
		{Type: token.INT, Literal: "4"},
		{Type: token.DIV, Literal: "/"},
		{Type: token.INT, Literal: "5"},
		{Type: token.MOD, Literal: "%"},
		{Type: token.LPAREN, Literal: "("},
		{Type: token.INT, Literal: "6"},
		{Type: token.RPAREN, Literal: ")"},
		{Type: token.INT, Literal: "3"},
		{Type: token.INC, Literal: "++"},
		{Type: token.ADD, Literal: "+"},
		{Type: token.INT, Literal: "2"},
		{Type: token.DEC, Literal: "--"},
		{Type: token.RBRACES, Literal: "}}"},
		{Type: token.EOF, Literal: ""},
	})
}

func TestStrings(t *testing.T) {
	inp := fmt.Sprintf(`{{ "Anna \"and\" Serhii" + %s }}`, "``")

	TokenizeString(t, inp, []token.Token{
		{Type: token.LBRACES, Literal: "{{"},
		{Type: token.STR, Literal: `Anna "and" Serhii`},
		{Type: token.ADD, Literal: "+"},
		{Type: token.STR, Literal: ""},
		{Type: token.RBRACES, Literal: "}}"},
		{Type: token.EOF, Literal: ""},
	})
}

func TestTernary(t *testing.T) {
	inp := `<small>{{ true ? "Yes!" : "No!" }}</small>`

	TokenizeString(t, inp, []token.Token{
		{Type: token.HTML, Literal: "<small>"},
		{Type: token.LBRACES, Literal: "{{"},
		{Type: token.TRUE, Literal: "true"},
		{Type: token.QUESTION, Literal: "?"},
		{Type: token.STR, Literal: "Yes!"},
		{Type: token.COLON, Literal: ":"},
		{Type: token.STR, Literal: "No!"},
		{Type: token.RBRACES, Literal: "}}"},
		{Type: token.HTML, Literal: "</small>"},
		{Type: token.EOF, Literal: ""},
	})
}

func TestIllegalToken(t *testing.T) {
	inp := `{{ 4 }`

	TokenizeString(t, inp, []token.Token{
		{Type: token.LBRACES, Literal: "{{"},
		{Type: token.INT, Literal: "4"},
		{Type: token.ILLEGAL, Literal: "}"},
	})
}

func TestVariableDeclaration(t *testing.T) {
	t.Run("Without var", func(tt *testing.T) {
		inp := `{{ a := 1 }}`

		TokenizeString(tt, inp, []token.Token{
			{Type: token.LBRACES, Literal: "{{"},
			{Type: token.IDENT, Literal: "a"},
			{Type: token.DEFINE, Literal: ":="},
			{Type: token.INT, Literal: "1"},
			{Type: token.RBRACES, Literal: "}}"},
			{Type: token.EOF, Literal: ""},
		})
	})

	t.Run("With var", func(tt *testing.T) {
		inp := `{{ var b = "hello" }}`

		TokenizeString(tt, inp, []token.Token{
			{Type: token.LBRACES, Literal: "{{"},
			{Type: token.VAR, Literal: "var"},
			{Type: token.IDENT, Literal: "b"},
			{Type: token.ASSIGN, Literal: "="},
			{Type: token.STR, Literal: "hello"},
			{Type: token.RBRACES, Literal: "}}"},
			{Type: token.EOF, Literal: ""},
		})
	})
}

func TestOther(t *testing.T) {
	inp := "{{ , == != <= >= > < }}"

	TokenizeString(t, inp, []token.Token{
		{Type: token.LBRACES, Literal: "{{"},
		{Type: token.COMMA, Literal: ","},
		{Type: token.EQ, Literal: "=="},
		{Type: token.NOT_EQ, Literal: "!="},
		{Type: token.LTHAN_EQ, Literal: "<="},
		{Type: token.GTHAN_EQ, Literal: ">="},
		{Type: token.GTHAN, Literal: ">"},
		{Type: token.LTHAN, Literal: "<"},
		{Type: token.RBRACES, Literal: "}}"},
		{Type: token.EOF, Literal: ""},
	})
}

func TestArray(t *testing.T) {
	inp := `{{ ["one", "two", "three"][1] }}`

	TokenizeString(t, inp, []token.Token{
		{Type: token.LBRACES, Literal: "{{"},
		{Type: token.LBRACKET, Literal: "["},
		{Type: token.STR, Literal: "one"},
		{Type: token.COMMA, Literal: ","},
		{Type: token.STR, Literal: "two"},
		{Type: token.COMMA, Literal: ","},
		{Type: token.STR, Literal: "three"},
		{Type: token.RBRACKET, Literal: "]"},
		{Type: token.LBRACKET, Literal: "["},
		{Type: token.INT, Literal: "1"},
		{Type: token.RBRACKET, Literal: "]"},
		{Type: token.RBRACES, Literal: "}}"},
		{Type: token.EOF, Literal: ""},
	})
}
