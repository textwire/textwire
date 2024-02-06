package lexer

import (
	"testing"

	"github.com/textwire/textwire/token"
)

func TokenizeString(t *testing.T, input string, expectTokens []token.Token) {
	l := New(input)

	for i, expectTok := range expectTokens {
		tok := l.NextToken()

		if tok.Literal != expectTok.Literal {
			t.Fatalf("tests[%d] - literal wrong. expected=\n'%s', got=\n'%s'",
				i, expectTok.Literal, tok.Literal)
		}

		if tok.Type != expectTok.Type {
			t.Fatalf("tests[%d] - tokentype wrong. expected='%q, got=%q",
				i, token.String(expectTok.Type), token.String(tok.Type))
		}
	}
}

func TestHTML(t *testing.T) {
	inp := `<h2 class="container">The winter is test@mail.com</h2>`

	TokenizeString(t, inp, []token.Token{
		{Type: token.HTML, Literal: `<h2 class="container">The winter is test@mail.com</h2>`},
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
	inp := "{{ testVar another_var12 nil false !true}}"

	TokenizeString(t, inp, []token.Token{
		{Type: token.LBRACES, Literal: "{{"},
		{Type: token.IDENT, Literal: "testVar"},
		{Type: token.IDENT, Literal: "another_var12"},
		{Type: token.NIL, Literal: "nil"},
		{Type: token.FALSE, Literal: "false"},
		{Type: token.NOT, Literal: "!"},
		{Type: token.TRUE, Literal: "true"},
		{Type: token.RBRACES, Literal: "}}"},
		{Type: token.EOF, Literal: ""},
	})
}

func TestIfStatement(t *testing.T) {
	inp := `@if(true(()))one@elseif(false){{ "nice" }}@elsethree@endfour`

	TokenizeString(t, inp, []token.Token{
		{Type: token.IF, Literal: "@if"},
		{Type: token.LPAREN, Literal: "("},
		{Type: token.TRUE, Literal: "true"},
		{Type: token.LPAREN, Literal: "("},
		{Type: token.LPAREN, Literal: "("},
		{Type: token.RPAREN, Literal: ")"},
		{Type: token.RPAREN, Literal: ")"},
		{Type: token.RPAREN, Literal: ")"},
		{Type: token.HTML, Literal: "one"},
		{Type: token.ELSEIF, Literal: "@elseif"},
		{Type: token.LPAREN, Literal: "("},
		{Type: token.FALSE, Literal: "false"},
		{Type: token.RPAREN, Literal: ")"},
		{Type: token.LBRACES, Literal: "{{"},
		{Type: token.STR, Literal: "nice"},
		{Type: token.RBRACES, Literal: "}}"},
		{Type: token.ELSE, Literal: "@else"},
		{Type: token.HTML, Literal: "three"},
		{Type: token.END, Literal: "@end"},
		{Type: token.HTML, Literal: "four"},
	})
}

func TestUseStatement(t *testing.T) {
	inp := `<div>@use("layouts/main")</div>`

	TokenizeString(t, inp, []token.Token{
		{Type: token.HTML, Literal: "<div>"},
		{Type: token.USE, Literal: "@use"},
		{Type: token.LPAREN, Literal: "("},
		{Type: token.STR, Literal: "layouts/main"},
		{Type: token.RPAREN, Literal: ")"},
		{Type: token.HTML, Literal: "</div>"},
		{Type: token.EOF, Literal: ""},
	})
}

func TestReserveStatement(t *testing.T) {
	inp := `<div>@reserve("title")</div>`

	TokenizeString(t, inp, []token.Token{
		{Type: token.HTML, Literal: "<div>"},
		{Type: token.RESERVE, Literal: "@reserve"},
		{Type: token.LPAREN, Literal: "("},
		{Type: token.STR, Literal: "title"},
		{Type: token.RPAREN, Literal: ")"},
		{Type: token.HTML, Literal: "</div>"},
		{Type: token.EOF, Literal: ""},
	})
}

func TestInsertStatement(t *testing.T) {
	inp := `@insert("title")<div>Nice one</div>@end`

	TokenizeString(t, inp, []token.Token{
		{Type: token.INSERT, Literal: "@insert"},
		{Type: token.LPAREN, Literal: "("},
		{Type: token.STR, Literal: "title"},
		{Type: token.RPAREN, Literal: ")"},
		{Type: token.HTML, Literal: "<div>Nice one</div>"},
		{Type: token.END, Literal: "@end"},
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
	inp := `{{ "Anna \"and\" Serhii" + '' }}`

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

func TestLineNumber(t *testing.T) {
	tests := []struct {
		inp  string
		line uint
	}{
		{"", 0},
		{" ", 1},
		{"\n", 2},
		{"1\n2\n3\n4", 4},
		{"{{ age := 3 }}", 1},
		{`{{ age := 3; age }}`, 1},
		{
			`<h1>Title</h1>
			<p>Test</p>`, 2,
		},
		{
			`<h1>Title</h1>
			<p>Test</p>
			{{ age := 3 }}
			{{ age }}`, 4,
		},
		{
			`{{
				age := 3;
				age
			}}`, 4,
		},
	}

	for _, tt := range tests {
		l := New(tt.inp)
		var lastTok token.Token

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			lastTok = tok
		}

		if lastTok.Line != tt.line {
			t.Errorf("Expected line number %d, got %d", tt.line, lastTok.Line)
		}
	}
}

func TestIsDirectoryStart(t *testing.T) {
	t.Run("Not a directive", func(tt *testing.T) {
		l := New(`test@email.com`)

		if ok := l.isDirectiveStmt(); ok {
			tt.Errorf("Expected not a directive")
		}
	})

	t.Run("Directive", func(tt *testing.T) {
		l := New(`@if(true)@end`)

		if ok := l.isDirectiveStmt(); !ok {
			t.Errorf("Expected a directive")
		}
	})
}

func TestFunctionCall(t *testing.T) {
	t.Run("On string", func(tt *testing.T) {
		inp := `{{ "test".upper() }}`

		TokenizeString(t, inp, []token.Token{
			{Type: token.LBRACES, Literal: "{{"},
			{Type: token.STR, Literal: "test"},
			{Type: token.DOT, Literal: "."},
			{Type: token.IDENT, Literal: "upper"},
			{Type: token.LPAREN, Literal: "("},
			{Type: token.RPAREN, Literal: ")"},
			{Type: token.RBRACES, Literal: "}}"},
			{Type: token.EOF, Literal: ""},
		})
	})

	t.Run("On float", func(tt *testing.T) {
		inp := `{{ 3.14.round() }}`

		TokenizeString(t, inp, []token.Token{
			{Type: token.LBRACES, Literal: "{{"},
			{Type: token.FLOAT, Literal: "3.14"},
			{Type: token.DOT, Literal: "."},
			{Type: token.IDENT, Literal: "round"},
			{Type: token.LPAREN, Literal: "("},
			{Type: token.RPAREN, Literal: ")"},
			{Type: token.RBRACES, Literal: "}}"},
			{Type: token.EOF, Literal: ""},
		})
	})
}
