package lexer

import (
	"testing"

	token "github.com/textwire/textwire/v2/token"
)

func TokenizeString(t *testing.T, input string, expectTokens []token.Token) {
	l := New(input)

	for i, expectTok := range expectTokens {
		tok := l.NextToken()

		if tok.Literal != expectTok.Literal {
			t.Fatalf("tests[%d] - literal wrong. expected='%s', got='%s'",
				i, expectTok.Literal, tok.Literal)
		}

		if tok.Type != expectTok.Type {
			t.Fatalf("tests[%d] - tokentype wrong. expected='%s', got='%s'",
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
	inp := "{{ 0.12 1.1111 9.1 }}"

	TokenizeString(t, inp, []token.Token{
		{Type: token.LBRACES, Literal: "{{"},
		{Type: token.FLOAT, Literal: "0.12"},
		{Type: token.FLOAT, Literal: "1.1111"},
		{Type: token.FLOAT, Literal: "9.1"},
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

func TestIfStmt(t *testing.T) {
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
		{Type: token.ELSE_IF, Literal: "@elseif"},
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

func TestUseStmt(t *testing.T) {
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

func TestReserveStmt(t *testing.T) {
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

func TestInsertStmt(t *testing.T) {
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

func TestStrings(test *testing.T) {
	test.Run("String with quotes", func(t *testing.T) {
		inp := `{{ "Anna \"and\" Serhii" + '' }}`

		TokenizeString(t, inp, []token.Token{
			{Type: token.LBRACES, Literal: "{{"},
			{Type: token.STR, Literal: `Anna "and" Serhii`},
			{Type: token.ADD, Literal: "+"},
			{Type: token.STR, Literal: ""},
			{Type: token.RBRACES, Literal: "}}"},
			{Type: token.EOF, Literal: ""},
		})
	})

	test.Run("Empty string", func(t *testing.T) {
		inp := `{{ "" }}`

		TokenizeString(t, inp, []token.Token{
			{Type: token.LBRACES, Literal: "{{"},
			{Type: token.STR, Literal: ""},
			{Type: token.RBRACES, Literal: "}}"},
			{Type: token.EOF, Literal: ""},
		})
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

func TestVariableDeclaration(t *testing.T) {
	inp := `{{ a = 1 }}`

	TokenizeString(t, inp, []token.Token{
		{Type: token.LBRACES, Literal: "{{"},
		{Type: token.IDENT, Literal: "a"},
		{Type: token.ASSIGN, Literal: "="},
		{Type: token.INT, Literal: "1"},
		{Type: token.RBRACES, Literal: "}}"},
		{Type: token.EOF, Literal: ""},
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

func TestDebugLineNumber(t *testing.T) {
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

	for _, tc := range tests {
		l := New(tc.inp)
		var lastTok token.Token

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			lastTok = tok
		}

		if lastTok.DebugLine != tc.line {
			t.Errorf("Expected line number %d, got %d", tc.line, lastTok.DebugLine)
		}
	}
}

func TestTokenPosition(t *testing.T) {
	tests := []struct {
		pos        token.Position
		tokenIndex int
	}{
		{
			token.Position{
				StartLine: 0,
				EndLine:   1,
				StartCol:  0,
				EndCol:    3,
			},
			0,
		},
	}

	inp := `<ul>
    {{ age = 3 }}
</ul>`

	for _, tc := range tests {
		l := New(inp)
		var targetTok token.Token

		for i := 0; i <= tc.tokenIndex; i++ {
			targetTok = l.NextToken()
		}

		if targetTok.Pos.StartCol != tc.pos.StartCol {
			t.Errorf("Expected token %s StartCol: %v, got %v",
				targetTok.Literal, tc.pos.StartCol, targetTok.Pos.StartCol)
		}

		if targetTok.Pos.EndCol != tc.pos.EndCol {
			t.Errorf("Expected token %s EndCol: %v, got %v",
				targetTok.Literal, tc.pos.EndCol, targetTok.Pos.EndCol)
		}

		if targetTok.Pos.EndLine != tc.pos.EndLine {
			t.Errorf("Expected token %s EndLine: %v, got %v",
				targetTok.Literal, tc.pos.EndLine, targetTok.Pos.EndLine)
		}

		if targetTok.Pos.StartLine != tc.pos.StartLine {
			t.Errorf("Expected token %s StartLine: %v, got %v",
				targetTok.Literal, tc.pos.StartLine, targetTok.Pos.StartLine)
		}
	}
}

func TestIsDirectoryStart(t *testing.T) {
	t.Run("Not a directive", func(t *testing.T) {
		l := New(`test@email.com`)

		if ok := l.isDirectiveStmt(); ok {
			t.Errorf("Expected not a directive")
		}
	})

	t.Run("Directive", func(t *testing.T) {
		l := New(`@if(true)@end`)

		if ok := l.isDirectiveStmt(); !ok {
			t.Errorf("Expected a directive")
		}
	})
}

func TestCallExp(t *testing.T) {
	t.Run("On string", func(t *testing.T) {
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

	t.Run("On int", func(t *testing.T) {
		inp := `{{ 3.int() }}`

		TokenizeString(t, inp, []token.Token{
			{Type: token.LBRACES, Literal: "{{"},
			{Type: token.INT, Literal: "3"},
			{Type: token.DOT, Literal: "."},
			{Type: token.IDENT, Literal: "int"},
			{Type: token.LPAREN, Literal: "("},
			{Type: token.RPAREN, Literal: ")"},
			{Type: token.RBRACES, Literal: "}}"},
			{Type: token.EOF, Literal: ""},
		})
	})
}

func TestForLoopStatement(t *testing.T) {
	inp := `@for(i = 0; i < 10; i++)`

	TokenizeString(t, inp, []token.Token{
		{Type: token.FOR, Literal: "@for"},
		{Type: token.LPAREN, Literal: "("},
		{Type: token.IDENT, Literal: "i"},
		{Type: token.ASSIGN, Literal: "="},
		{Type: token.INT, Literal: "0"},
		{Type: token.SEMI, Literal: ";"},
		{Type: token.IDENT, Literal: "i"},
		{Type: token.LTHAN, Literal: "<"},
		{Type: token.INT, Literal: "10"},
		{Type: token.SEMI, Literal: ";"},
		{Type: token.IDENT, Literal: "i"},
		{Type: token.INC, Literal: "++"},
		{Type: token.RPAREN, Literal: ")"},
		{Type: token.EOF, Literal: ""},
	})
}

func TestEachLoopStatement(t *testing.T) {
	inp := `@each(n in [1, 2, 3])@if(n == 2)@break@end{{ n }}@end`

	TokenizeString(t, inp, []token.Token{
		{Type: token.EACH, Literal: "@each"},
		{Type: token.LPAREN, Literal: "("},
		{Type: token.IDENT, Literal: "n"},
		{Type: token.IN, Literal: "in"},
		{Type: token.LBRACKET, Literal: "["},
		{Type: token.INT, Literal: "1"},
		{Type: token.COMMA, Literal: ","},
		{Type: token.INT, Literal: "2"},
		{Type: token.COMMA, Literal: ","},
		{Type: token.INT, Literal: "3"},
		{Type: token.RBRACKET, Literal: "]"},
		{Type: token.RPAREN, Literal: ")"},
		{Type: token.IF, Literal: "@if"},
		{Type: token.LPAREN, Literal: "("},
		{Type: token.IDENT, Literal: "n"},
		{Type: token.EQ, Literal: "=="},
		{Type: token.INT, Literal: "2"},
		{Type: token.RPAREN, Literal: ")"},
		{Type: token.BREAK, Literal: "@break"},
		{Type: token.END, Literal: "@end"},
		{Type: token.LBRACES, Literal: "{{"},
		{Type: token.IDENT, Literal: "n"},
		{Type: token.RBRACES, Literal: "}}"},
		{Type: token.END, Literal: "@end"},
		{Type: token.EOF, Literal: ""},
	})
}

func TestObjectStatement(t *testing.T) {
	inp := `{{ {"father": {"name": "John"}} }}`

	TokenizeString(t, inp, []token.Token{
		{Type: token.LBRACES, Literal: "{{"},
		{Type: token.LBRACE, Literal: "{"},
		{Type: token.STR, Literal: "father"},
		{Type: token.COLON, Literal: ":"},
		{Type: token.LBRACE, Literal: "{"},
		{Type: token.STR, Literal: "name"},
		{Type: token.COLON, Literal: ":"},
		{Type: token.STR, Literal: "John"},
		{Type: token.RBRACE, Literal: "}"},
		{Type: token.RBRACE, Literal: "}"},
		{Type: token.RBRACES, Literal: "}}"},
		{Type: token.EOF, Literal: ""},
	})
}

func TestBreakDirectives(t *testing.T) {
	inp := `@breakIf(true) @break @continue @continueIf(false)`

	TokenizeString(t, inp, []token.Token{
		{Type: token.BREAK_IF, Literal: "@breakIf"},
		{Type: token.LPAREN, Literal: "("},
		{Type: token.TRUE, Literal: "true"},
		{Type: token.RPAREN, Literal: ")"},
		{Type: token.HTML, Literal: " "},
		{Type: token.BREAK, Literal: "@break"},
		{Type: token.HTML, Literal: " "},
		{Type: token.CONTINUE, Literal: "@continue"},
		{Type: token.HTML, Literal: " "},
		{Type: token.CONTINUE_IF, Literal: "@continueIf"},
		{Type: token.LPAREN, Literal: "("},
		{Type: token.FALSE, Literal: "false"},
		{Type: token.RPAREN, Literal: ")"},
		{Type: token.EOF, Literal: ""},
	})
}

func TestComponentDirective(t *testing.T) {
	inp := `@component("components/book-card", { c: card })`

	TokenizeString(t, inp, []token.Token{
		{Type: token.COMPONENT, Literal: "@component"},
		{Type: token.LPAREN, Literal: "("},
		{Type: token.STR, Literal: "components/book-card"},
		{Type: token.COMMA, Literal: ","},
		{Type: token.LBRACE, Literal: "{"},
		{Type: token.IDENT, Literal: "c"},
		{Type: token.COLON, Literal: ":"},
		{Type: token.IDENT, Literal: "card"},
		{Type: token.RBRACE, Literal: "}"},
		{Type: token.RPAREN, Literal: ")"},
		{Type: token.EOF, Literal: ""},
	})
}

func TestComponentSlotDirective(t *testing.T) {
	t.Run("slots with parentheses", func(t *testing.T) {
		inp := `@component("card")@slot("top")<h1>Hello</h1>@end@end`

		TokenizeString(t, inp, []token.Token{
			{Type: token.COMPONENT, Literal: "@component"},
			{Type: token.LPAREN, Literal: "("},
			{Type: token.STR, Literal: "card"},
			{Type: token.RPAREN, Literal: ")"},
			{Type: token.SLOT, Literal: "@slot"},
			{Type: token.LPAREN, Literal: "("},
			{Type: token.STR, Literal: "top"},
			{Type: token.RPAREN, Literal: ")"},
			{Type: token.HTML, Literal: "<h1>Hello</h1>"},
			{Type: token.END, Literal: "@end"},
			{Type: token.END, Literal: "@end"},
			{Type: token.EOF, Literal: ""},
		})
	})

	t.Run("slots without parentheses", func(t *testing.T) {
		inp := `@component("card")@slotNICE@end@end`

		TokenizeString(t, inp, []token.Token{
			{Type: token.COMPONENT, Literal: "@component"},
			{Type: token.LPAREN, Literal: "("},
			{Type: token.STR, Literal: "card"},
			{Type: token.RPAREN, Literal: ")"},
			{Type: token.SLOT, Literal: "@slot"},
			{Type: token.HTML, Literal: "NICE"},
			{Type: token.END, Literal: "@end"},
			{Type: token.END, Literal: "@end"},
			{Type: token.EOF, Literal: ""},
		})
	})
}

// Comments should be ignored by the lexer
func TestCommentStatement(t *testing.T) {
	t.Run("Simple comment", func(t *testing.T) {
		inp := `<div>{{-- This is a comment --}}</div>`

		TokenizeString(t, inp, []token.Token{
			{Type: token.HTML, Literal: "<div>"},
			{Type: token.HTML, Literal: "</div>"},
			{Type: token.EOF, Literal: ""},
		})
	})

	t.Run("Commented code", func(t *testing.T) {
		inp := `{{-- @each(u in users){{ u }}@end --}}`

		TokenizeString(t, inp, []token.Token{
			{Type: token.EOF, Literal: ""},
		})
	})
}
