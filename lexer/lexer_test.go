package lexer

import (
	"testing"

	token "github.com/textwire/textwire/v2/token"
)

func TokenizeString(t *testing.T, input string, expectTokens []token.Token) {
	l := New(input)

	for _, expectTok := range expectTokens {
		tok := l.NextToken()

		if tok.Literal != expectTok.Literal {
			t.Fatalf("token %q - literal wrong. expected='%s', got='%s'",
				tok.Literal, expectTok.Literal, tok.Literal)
		}

		if tok.Type != expectTok.Type {
			t.Fatalf("token %q - tokentype wrong. expected='%s', got='%s'",
				tok.Literal, token.String(expectTok.Type), token.String(tok.Type))
		}

		if tok.Pos != expectTok.Pos {
			t.Fatalf("token %q - position wrong.\nexpected: {startCol=%d, endCol=%d, startLine=%d, endLine=%d}\ngot:      {startCol=%d, endCol=%d, startLine=%d, endLine=%d}",
				tok.Literal, expectTok.Pos.StartCol, expectTok.Pos.EndCol,
				expectTok.Pos.StartLine, expectTok.Pos.EndLine,
				tok.Pos.StartCol, tok.Pos.EndCol, tok.Pos.StartLine,
				tok.Pos.EndLine)
		}
	}
}

func TestHTML(t *testing.T) {
	inp := `<h2 class="container">The winter is test@mail.com</h2>`

	TokenizeString(t, inp, []token.Token{
		{
			Type:    token.HTML,
			Literal: `<h2 class="container">The winter is test@mail.com</h2>`,
			Pos:     token.Position{EndCol: 53},
		},
		{
			Type:    token.EOF,
			Literal: "",
			Pos:     token.Position{StartCol: 54, EndCol: 54},
		},
	})
}

func TestIntegers(t *testing.T) {
	inp := "<div>{{ 0 1 2 3 4 5 6 7 8 9 234 -41 }}</div>"

	TokenizeString(t, inp, []token.Token{
		{Type: token.HTML, Literal: "<div>", Pos: token.Position{EndCol: 4}},
		{Type: token.LBRACES, Literal: "{{", Pos: token.Position{StartCol: 5, EndCol: 6}},
		{Type: token.INT, Literal: "0", Pos: token.Position{StartCol: 8, EndCol: 8}},
		{Type: token.INT, Literal: "1", Pos: token.Position{StartCol: 10, EndCol: 10}},
		{Type: token.INT, Literal: "2", Pos: token.Position{StartCol: 12, EndCol: 12}},
		{Type: token.INT, Literal: "3", Pos: token.Position{StartCol: 14, EndCol: 14}},
		{Type: token.INT, Literal: "4", Pos: token.Position{StartCol: 16, EndCol: 16}},
		{Type: token.INT, Literal: "5", Pos: token.Position{StartCol: 18, EndCol: 18}},
		{Type: token.INT, Literal: "6", Pos: token.Position{StartCol: 20, EndCol: 20}},
		{Type: token.INT, Literal: "7", Pos: token.Position{StartCol: 22, EndCol: 22}},
		{Type: token.INT, Literal: "8", Pos: token.Position{StartCol: 24, EndCol: 24}},
		{Type: token.INT, Literal: "9", Pos: token.Position{StartCol: 26, EndCol: 26}},
		{Type: token.INT, Literal: "234", Pos: token.Position{StartCol: 28, EndCol: 30}},
		{Type: token.SUB, Literal: "-", Pos: token.Position{StartCol: 32, EndCol: 32}},
		{Type: token.INT, Literal: "41", Pos: token.Position{StartCol: 33, EndCol: 34}},
		{Type: token.RBRACES, Literal: "}}", Pos: token.Position{StartCol: 36, EndCol: 37}},
		{Type: token.HTML, Literal: "</div>", Pos: token.Position{StartCol: 38, EndCol: 43}},
		{Type: token.EOF, Literal: "", Pos: token.Position{StartCol: 44, EndCol: 44}},
	})
}

func TestFloats(t *testing.T) {
	inp := "{{ 0.12 1.1111 9.1 }}"

	TokenizeString(t, inp, []token.Token{
		{Type: token.LBRACES, Literal: "{{", Pos: token.Position{EndCol: 1}},
		{Type: token.FLOAT, Literal: "0.12", Pos: token.Position{StartCol: 3, EndCol: 6}},
		{Type: token.FLOAT, Literal: "1.1111", Pos: token.Position{StartCol: 8, EndCol: 13}},
		{Type: token.FLOAT, Literal: "9.1", Pos: token.Position{StartCol: 15, EndCol: 17}},
		{Type: token.RBRACES, Literal: "}}", Pos: token.Position{StartCol: 19, EndCol: 20}},
		{Type: token.EOF, Literal: "", Pos: token.Position{StartCol: 21, EndCol: 21}},
	})
}

func TestIdentifiers(t *testing.T) {
	inp := "{{ testVar another_var12 nil false !true}}"

	TokenizeString(t, inp, []token.Token{
		{Type: token.LBRACES, Literal: "{{", Pos: token.Position{EndCol: 1}},
		{Type: token.IDENT, Literal: "testVar", Pos: token.Position{StartCol: 3, EndCol: 9}},
		{Type: token.IDENT, Literal: "another_var12", Pos: token.Position{StartCol: 11, EndCol: 23}},
		{Type: token.NIL, Literal: "nil", Pos: token.Position{StartCol: 25, EndCol: 27}},
		{Type: token.FALSE, Literal: "false", Pos: token.Position{StartCol: 29, EndCol: 33}},
		{Type: token.NOT, Literal: "!", Pos: token.Position{StartCol: 35, EndCol: 35}},
		{Type: token.TRUE, Literal: "true", Pos: token.Position{StartCol: 36, EndCol: 39}},
		{Type: token.RBRACES, Literal: "}}", Pos: token.Position{StartCol: 40, EndCol: 41}},
		{Type: token.EOF, Literal: "", Pos: token.Position{StartCol: 42, EndCol: 42}},
	})
}

func TestIfStmt(t *testing.T) {
	inp := `@if(true(()))one@elseif(false){{ "nice" }}@elsethree@endfour`

	TokenizeString(t, inp, []token.Token{
		{Type: token.IF, Literal: "@if", Pos: token.Position{EndCol: 2}},
		{Type: token.LPAREN, Literal: "(", Pos: token.Position{StartCol: 3, EndCol: 3}},
		{Type: token.TRUE, Literal: "true", Pos: token.Position{StartCol: 4, EndCol: 7}},
		{Type: token.LPAREN, Literal: "(", Pos: token.Position{StartCol: 8, EndCol: 8}},
		{Type: token.LPAREN, Literal: "(", Pos: token.Position{StartCol: 9, EndCol: 9}},
		{Type: token.RPAREN, Literal: ")", Pos: token.Position{StartCol: 10, EndCol: 10}},
		{Type: token.RPAREN, Literal: ")", Pos: token.Position{StartCol: 11, EndCol: 11}},
		{Type: token.RPAREN, Literal: ")", Pos: token.Position{StartCol: 12, EndCol: 12}},
		{Type: token.HTML, Literal: "one", Pos: token.Position{StartCol: 13, EndCol: 15}},
		{Type: token.ELSE_IF, Literal: "@elseif", Pos: token.Position{StartCol: 16, EndCol: 22}},
		{Type: token.LPAREN, Literal: "(", Pos: token.Position{StartCol: 23, EndCol: 23}},
		{Type: token.FALSE, Literal: "false", Pos: token.Position{StartCol: 24, EndCol: 28}},
		{Type: token.RPAREN, Literal: ")", Pos: token.Position{StartCol: 29, EndCol: 29}},
		{Type: token.LBRACES, Literal: "{{", Pos: token.Position{StartCol: 30, EndCol: 31}},
		{Type: token.STR, Literal: "nice", Pos: token.Position{StartCol: 33, EndCol: 38}},
		{Type: token.RBRACES, Literal: "}}", Pos: token.Position{StartCol: 40, EndCol: 41}},
		{Type: token.ELSE, Literal: "@else", Pos: token.Position{StartCol: 42, EndCol: 46}},
		{Type: token.HTML, Literal: "three", Pos: token.Position{StartCol: 47, EndCol: 51}},
		{Type: token.END, Literal: "@end", Pos: token.Position{StartCol: 52, EndCol: 55}},
		{Type: token.HTML, Literal: "four", Pos: token.Position{StartCol: 56, EndCol: 59}},
		{Type: token.EOF, Literal: "", Pos: token.Position{StartCol: 60, EndCol: 60}},
	})
}

func TestUseStmt(t *testing.T) {
	inp := `<div>@use("layouts/main")</div>`

	TokenizeString(t, inp, []token.Token{
		{Type: token.HTML, Literal: "<div>", Pos: token.Position{EndCol: 4}},
		{Type: token.USE, Literal: "@use", Pos: token.Position{StartCol: 5, EndCol: 8}},
		{Type: token.LPAREN, Literal: "(", Pos: token.Position{StartCol: 9, EndCol: 9}},
		{Type: token.STR, Literal: "layouts/main", Pos: token.Position{StartCol: 10, EndCol: 23}},
		{Type: token.RPAREN, Literal: ")", Pos: token.Position{StartCol: 24, EndCol: 24}},
		{Type: token.HTML, Literal: "</div>", Pos: token.Position{StartCol: 25, EndCol: 30}},
		{Type: token.EOF, Literal: "", Pos: token.Position{StartCol: 31, EndCol: 31}},
	})
}

func TestReserveStmt(t *testing.T) {
	inp := `<div>@reserve("title")</div>`

	TokenizeString(t, inp, []token.Token{
		{Type: token.HTML, Literal: "<div>", Pos: token.Position{EndCol: 4}},
		{Type: token.RESERVE, Literal: "@reserve", Pos: token.Position{StartCol: 5, EndCol: 12}},
		{Type: token.LPAREN, Literal: "(", Pos: token.Position{StartCol: 13, EndCol: 13}},
		{Type: token.STR, Literal: "title", Pos: token.Position{StartCol: 14, EndCol: 20}},
		{Type: token.RPAREN, Literal: ")", Pos: token.Position{StartCol: 21, EndCol: 21}},
		{Type: token.HTML, Literal: "</div>", Pos: token.Position{StartCol: 22, EndCol: 27}},
		{Type: token.EOF, Literal: "", Pos: token.Position{StartCol: 28, EndCol: 28}},
	})
}

func TestInsertStmt(t *testing.T) {
	inp := `@insert("title")<div>Nice one</div>@end`

	TokenizeString(t, inp, []token.Token{
		{Type: token.INSERT, Literal: "@insert", Pos: token.Position{EndCol: 6}},
		{Type: token.LPAREN, Literal: "(", Pos: token.Position{StartCol: 7, EndCol: 7}},
		{Type: token.STR, Literal: "title", Pos: token.Position{StartCol: 8, EndCol: 14}},
		{Type: token.RPAREN, Literal: ")", Pos: token.Position{StartCol: 15, EndCol: 15}},
		{Type: token.HTML, Literal: "<div>Nice one</div>", Pos: token.Position{StartCol: 16, EndCol: 34}},
		{Type: token.END, Literal: "@end", Pos: token.Position{StartCol: 35, EndCol: 38}},
		{Type: token.EOF, Literal: "", Pos: token.Position{StartCol: 39, EndCol: 39}},
	})
}

func TestOperators(t *testing.T) {
	inp := "{{ 1 + 2 - 3 * 4 / 5 % (6) 3++ + 2-- }}"

	TokenizeString(t, inp, []token.Token{
		{Type: token.LBRACES, Literal: "{{", Pos: token.Position{EndCol: 1}},
		{Type: token.INT, Literal: "1", Pos: token.Position{StartCol: 3, EndCol: 3}},
		{Type: token.ADD, Literal: "+", Pos: token.Position{StartCol: 5, EndCol: 5}},
		{Type: token.INT, Literal: "2", Pos: token.Position{StartCol: 7, EndCol: 7}},
		{Type: token.SUB, Literal: "-", Pos: token.Position{StartCol: 9, EndCol: 9}},
		{Type: token.INT, Literal: "3", Pos: token.Position{StartCol: 11, EndCol: 11}},
		{Type: token.MUL, Literal: "*", Pos: token.Position{StartCol: 13, EndCol: 13}},
		{Type: token.INT, Literal: "4", Pos: token.Position{StartCol: 15, EndCol: 15}},
		{Type: token.DIV, Literal: "/", Pos: token.Position{StartCol: 17, EndCol: 17}},
		{Type: token.INT, Literal: "5", Pos: token.Position{StartCol: 19, EndCol: 19}},
		{Type: token.MOD, Literal: "%", Pos: token.Position{StartCol: 21, EndCol: 21}},
		{Type: token.LPAREN, Literal: "(", Pos: token.Position{StartCol: 23, EndCol: 23}},
		{Type: token.INT, Literal: "6", Pos: token.Position{StartCol: 24, EndCol: 24}},
		{Type: token.RPAREN, Literal: ")", Pos: token.Position{StartCol: 25, EndCol: 25}},
		{Type: token.INT, Literal: "3", Pos: token.Position{StartCol: 27, EndCol: 27}},
		{Type: token.INC, Literal: "++", Pos: token.Position{StartCol: 28, EndCol: 29}},
		{Type: token.ADD, Literal: "+", Pos: token.Position{StartCol: 31, EndCol: 31}},
		{Type: token.INT, Literal: "2", Pos: token.Position{StartCol: 33, EndCol: 33}},
		{Type: token.DEC, Literal: "--", Pos: token.Position{StartCol: 34, EndCol: 35}},
		{Type: token.RBRACES, Literal: "}}", Pos: token.Position{StartCol: 37, EndCol: 38}},
		{Type: token.EOF, Literal: "", Pos: token.Position{StartCol: 39, EndCol: 39}},
	})
}

func TestStrings(test *testing.T) {
	test.Run("String with quotes", func(t *testing.T) {
		inp := `{{ "Anna \"and\" Serhii" + '' }}`

		TokenizeString(t, inp, []token.Token{
			{Type: token.LBRACES, Literal: "{{", Pos: token.Position{EndCol: 1}},
			{Type: token.STR, Literal: `Anna "and" Serhii`, Pos: token.Position{StartCol: 3, EndCol: 23}},
			{Type: token.ADD, Literal: "+", Pos: token.Position{StartCol: 25, EndCol: 25}},
			{Type: token.STR, Literal: "", Pos: token.Position{StartCol: 27, EndCol: 28}},
			{Type: token.RBRACES, Literal: "}}", Pos: token.Position{StartCol: 30, EndCol: 31}},
			{Type: token.EOF, Literal: "", Pos: token.Position{StartCol: 32, EndCol: 32}},
		})
	})

	test.Run("String reads correctly with braces", func(t *testing.T) {
		inp := `{{ "\{ {" }}`

		TokenizeString(t, inp, []token.Token{
			{Type: token.LBRACES, Literal: "{{", Pos: token.Position{EndCol: 1}},
			{Type: token.STR, Literal: `\{ {`, Pos: token.Position{StartCol: 3, EndCol: 8}},
			{Type: token.RBRACES, Literal: "}}", Pos: token.Position{StartCol: 10, EndCol: 11}},
			{Type: token.EOF, Literal: "", Pos: token.Position{StartCol: 12, EndCol: 12}},
		})
	})

	test.Run("Empty string", func(t *testing.T) {
		inp := `{{ "" }}`

		TokenizeString(t, inp, []token.Token{
			{Type: token.LBRACES, Literal: "{{", Pos: token.Position{EndCol: 1}},
			{Type: token.STR, Literal: "", Pos: token.Position{StartCol: 3, EndCol: 4}},
			{Type: token.RBRACES, Literal: "}}", Pos: token.Position{StartCol: 6, EndCol: 7}},
			{Type: token.EOF, Literal: "", Pos: token.Position{StartCol: 8, EndCol: 8}},
		})
	})
}

func TestTernary(t *testing.T) {
	inp := `<small>{{ true ? "Yes!" : "No!" }}</small>`

	TokenizeString(t, inp, []token.Token{
		{Type: token.HTML, Literal: "<small>", Pos: token.Position{EndCol: 6}},
		{Type: token.LBRACES, Literal: "{{", Pos: token.Position{StartCol: 7, EndCol: 8}},
		{Type: token.TRUE, Literal: "true", Pos: token.Position{StartCol: 10, EndCol: 13}},
		{Type: token.QUESTION, Literal: "?", Pos: token.Position{StartCol: 15, EndCol: 15}},
		{Type: token.STR, Literal: "Yes!", Pos: token.Position{StartCol: 17, EndCol: 22}},
		{Type: token.COLON, Literal: ":", Pos: token.Position{StartCol: 24, EndCol: 24}},
		{Type: token.STR, Literal: "No!", Pos: token.Position{StartCol: 26, EndCol: 30}},
		{Type: token.RBRACES, Literal: "}}", Pos: token.Position{StartCol: 32, EndCol: 33}},
		{Type: token.HTML, Literal: "</small>", Pos: token.Position{StartCol: 34, EndCol: 41}},
		{Type: token.EOF, Literal: "", Pos: token.Position{StartCol: 42, EndCol: 42}},
	})
}

func TestVariableDeclaration(t *testing.T) {
	inp := `{{ a = 1 }}`

	TokenizeString(t, inp, []token.Token{
		{Type: token.LBRACES, Literal: "{{", Pos: token.Position{EndCol: 1}},
		{Type: token.IDENT, Literal: "a", Pos: token.Position{StartCol: 3, EndCol: 3}},
		{Type: token.ASSIGN, Literal: "=", Pos: token.Position{StartCol: 5, EndCol: 5}},
		{Type: token.INT, Literal: "1", Pos: token.Position{StartCol: 7, EndCol: 7}},
		{Type: token.RBRACES, Literal: "}}", Pos: token.Position{StartCol: 9, EndCol: 10}},
		{Type: token.EOF, Literal: "", Pos: token.Position{StartCol: 11, EndCol: 11}},
	})
}

func TestLogicalAndOperator(t *testing.T) {
	inp := `{{ true && false }}`

	TokenizeString(t, inp, []token.Token{
		{Type: token.LBRACES, Literal: "{{", Pos: token.Position{EndCol: 1}},
		{Type: token.TRUE, Literal: "true", Pos: token.Position{StartCol: 3, EndCol: 6}},
		{Type: token.AND, Literal: "&&", Pos: token.Position{StartCol: 8, EndCol: 9}},
		{Type: token.FALSE, Literal: "false", Pos: token.Position{StartCol: 11, EndCol: 15}},
		{Type: token.RBRACES, Literal: "}}", Pos: token.Position{StartCol: 17, EndCol: 18}},
		{Type: token.EOF, Literal: "", Pos: token.Position{StartCol: 19, EndCol: 19}},
	})
}

func TestLogicalOrOperator(t *testing.T) {
	inp := `{{ true || false }}`

	TokenizeString(t, inp, []token.Token{
		{Type: token.LBRACES, Literal: "{{", Pos: token.Position{EndCol: 1}},
		{Type: token.TRUE, Literal: "true", Pos: token.Position{StartCol: 3, EndCol: 6}},
		{Type: token.OR, Literal: "||", Pos: token.Position{StartCol: 8, EndCol: 9}},
		{Type: token.FALSE, Literal: "false", Pos: token.Position{StartCol: 11, EndCol: 15}},
		{Type: token.RBRACES, Literal: "}}", Pos: token.Position{StartCol: 17, EndCol: 18}},
		{Type: token.EOF, Literal: "", Pos: token.Position{StartCol: 19, EndCol: 19}},
	})
}

func TestOther(t *testing.T) {
	inp := "{{ , == != <= >= > < }}"

	TokenizeString(t, inp, []token.Token{
		{Type: token.LBRACES, Literal: "{{", Pos: token.Position{EndCol: 1}},
		{Type: token.COMMA, Literal: ",", Pos: token.Position{StartCol: 3, EndCol: 3}},
		{Type: token.EQ, Literal: "==", Pos: token.Position{StartCol: 5, EndCol: 6}},
		{Type: token.NOT_EQ, Literal: "!=", Pos: token.Position{StartCol: 8, EndCol: 9}},
		{Type: token.LTHAN_EQ, Literal: "<=", Pos: token.Position{StartCol: 11, EndCol: 12}},
		{Type: token.GTHAN_EQ, Literal: ">=", Pos: token.Position{StartCol: 14, EndCol: 15}},
		{Type: token.GTHAN, Literal: ">", Pos: token.Position{StartCol: 17, EndCol: 17}},
		{Type: token.LTHAN, Literal: "<", Pos: token.Position{StartCol: 19, EndCol: 19}},
		{Type: token.RBRACES, Literal: "}}", Pos: token.Position{StartCol: 21, EndCol: 22}},
		{Type: token.EOF, Literal: "", Pos: token.Position{StartCol: 23, EndCol: 23}},
	})
}

func TestArray(t *testing.T) {
	inp := `{{ ["one", "two", "three"][1] }}`

	TokenizeString(t, inp, []token.Token{
		{Type: token.LBRACES, Literal: "{{", Pos: token.Position{EndCol: 1}},
		{Type: token.LBRACKET, Literal: "[", Pos: token.Position{StartCol: 3, EndCol: 3}},
		{Type: token.STR, Literal: "one", Pos: token.Position{StartCol: 4, EndCol: 8}},
		{Type: token.COMMA, Literal: ",", Pos: token.Position{StartCol: 9, EndCol: 9}},
		{Type: token.STR, Literal: "two", Pos: token.Position{StartCol: 11, EndCol: 15}},
		{Type: token.COMMA, Literal: ",", Pos: token.Position{StartCol: 16, EndCol: 16}},
		{Type: token.STR, Literal: "three", Pos: token.Position{StartCol: 18, EndCol: 24}},
		{Type: token.RBRACKET, Literal: "]", Pos: token.Position{StartCol: 25, EndCol: 25}},
		{Type: token.LBRACKET, Literal: "[", Pos: token.Position{StartCol: 26, EndCol: 26}},
		{Type: token.INT, Literal: "1", Pos: token.Position{StartCol: 27, EndCol: 27}},
		{Type: token.RBRACKET, Literal: "]", Pos: token.Position{StartCol: 28, EndCol: 28}},
		{Type: token.RBRACES, Literal: "}}", Pos: token.Position{StartCol: 30, EndCol: 31}},
		{Type: token.EOF, Literal: "", Pos: token.Position{StartCol: 32, EndCol: 32}},
	})
}

func TestErrorLineNumber(t *testing.T) {
	tests := []struct {
		inp  string
		line uint
	}{
		{"", 1},
		{" ", 1},
		{"\n", 1},
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

		if lastTok.ErrorLine() != tc.line {
			t.Errorf("Expected line number %d, got %d", tc.line, lastTok.ErrorLine())
		}
	}
}

func TestTokenPosition(t *testing.T) {
	tests := []struct {
		startL uint
		endL   uint
		startC uint
		endC   uint
	}{
		{startL: 0, endL: 1, startC: 0, endC: 3},   // <div>\n____
		{startL: 1, endL: 1, startC: 4, endC: 5},   // {{
		{startL: 1, endL: 1, startC: 7, endC: 9},   // age
		{startL: 1, endL: 1, startC: 11, endC: 11}, // =
		{startL: 1, endL: 1, startC: 13, endC: 18}, // 323.24
		{startL: 1, endL: 1, startC: 20, endC: 21}, // }}
		{startL: 1, endL: 2, startC: 22, endC: 3},  // \n____
		{startL: 2, endL: 2, startC: 4, endC: 11},  // @reserve
		{startL: 2, endL: 2, startC: 12, endC: 12}, // (
		{startL: 2, endL: 2, startC: 13, endC: 19}, // "title"
		{startL: 2, endL: 2, startC: 20, endC: 20}, // )
		{startL: 2, endL: 3, startC: 21, endC: 6},  // \n</div>\n
		{startL: 4, endL: 4, startC: 0, endC: 5},   // @break
	}

	inp := `<div>
    {{ age = 323.24 }}
    @reserve("title")
</div>
@break`

	for tokenIdx, tc := range tests {
		l := New(inp)
		var targetTok token.Token

		for i := 0; i <= tokenIdx; i++ {
			targetTok = l.NextToken()
		}

		pos := token.Position{
			StartLine: tc.startL,
			EndLine:   tc.endL,
			StartCol:  tc.startC,
			EndCol:    tc.endC,
		}

		if targetTok.Pos.StartCol != pos.StartCol {
			t.Errorf("Expected token %q StartCol: %d, got %d",
				targetTok.Literal, pos.StartCol, targetTok.Pos.StartCol)
		}

		if targetTok.Pos.EndCol != pos.EndCol {
			t.Errorf("Expected token %q EndCol: %d, got: %d",
				targetTok.Literal, pos.EndCol, targetTok.Pos.EndCol)
		}

		if targetTok.Pos.EndLine != pos.EndLine {
			t.Errorf("Expected token %q EndLine: %d, got: %d",
				targetTok.Literal, pos.EndLine, targetTok.Pos.EndLine)
		}

		if targetTok.Pos.StartLine != pos.StartLine {
			t.Errorf("Expected token %q StartLine: %d, got %d",
				targetTok.Literal, pos.StartLine, targetTok.Pos.StartLine)
		}
	}
}

func TestIsDirectoryToken(t *testing.T) {
	t.Run("Not a directive token", func(t *testing.T) {
		input := "test@email.com"
		l := New(input)

		isDir, escaped := l.isDirectiveToken()
		if isDir {
			t.Errorf("Expected %q not to be a directive token", input)
		}

		if escaped {
			t.Errorf("Expected %q not to be escaped directive token", input)
		}
	})

	t.Run("Directive token", func(t *testing.T) {
		input := "@break"
		l := New(input)

		isDir, escaped := l.isDirectiveToken()
		if !isDir {
			t.Errorf("Expected %q to be a directive token", input)
		}

		if escaped {
			t.Errorf("Expected %q not to be an escaped directive token", input)
		}
	})

	t.Run("Escaped directive token", func(t *testing.T) {
		input := `\@break`
		l := New(input)
		l.readChar() // skip the backslash

		isDir, escaped := l.isDirectiveToken()
		if isDir {
			t.Errorf("Expected %q not to be a directive", input)
		}

		if !escaped {
			t.Errorf("Expected %q to be escaped directive", input)
		}
	})
}

func Test_areBracesToken(t *testing.T) {
	t.Run("Not braces token", func(t *testing.T) {
		input := "some {{ text"
		l := New(input)

		areBraces, escaped := l.areBracesToken()
		if areBraces {
			t.Errorf("Expected %q not to be braces", input)
		}

		if escaped {
			t.Errorf("Expected %q not to be escaped braces", input)
		}
	})

	t.Run("Braces token", func(t *testing.T) {
		input := "{{ 123 }}"
		l := New(input)

		areBraces, escaped := l.areBracesToken()
		if !areBraces {
			t.Errorf("Expected %q to be braces token", input)
		}

		if escaped {
			t.Errorf("Expected %q not to be escaped braces token", input)
		}
	})

	t.Run("Escaped braces token", func(t *testing.T) {
		input := `\{{ 234 }}`
		l := New(input)
		l.readChar() // skip the backslash

		areBraces, escaped := l.areBracesToken()
		if areBraces {
			t.Errorf("Expected %q not to be braces token", input)
		}

		if !escaped {
			t.Errorf("Expected %q to be escaped braces token", input)
		}
	})
}

func TestCallExp(t *testing.T) {
	t.Run("On string", func(t *testing.T) {
		inp := `{{ "test".upper() }}`

		TokenizeString(t, inp, []token.Token{
			{Type: token.LBRACES, Literal: "{{", Pos: token.Position{EndCol: 1}},
			{Type: token.STR, Literal: "test", Pos: token.Position{StartCol: 3, EndCol: 8}},
			{Type: token.DOT, Literal: ".", Pos: token.Position{StartCol: 9, EndCol: 9}},
			{Type: token.IDENT, Literal: "upper", Pos: token.Position{StartCol: 10, EndCol: 14}},
			{Type: token.LPAREN, Literal: "(", Pos: token.Position{StartCol: 15, EndCol: 15}},
			{Type: token.RPAREN, Literal: ")", Pos: token.Position{StartCol: 16, EndCol: 16}},
			{Type: token.RBRACES, Literal: "}}", Pos: token.Position{StartCol: 18, EndCol: 19}},
			{Type: token.EOF, Literal: "", Pos: token.Position{StartCol: 20, EndCol: 20}},
		})
	})

	t.Run("On int", func(t *testing.T) {
		inp := `{{ 3.int() }}`

		TokenizeString(t, inp, []token.Token{
			{Type: token.LBRACES, Literal: "{{", Pos: token.Position{EndCol: 1}},
			{Type: token.INT, Literal: "3", Pos: token.Position{StartCol: 3, EndCol: 3}},
			{Type: token.DOT, Literal: ".", Pos: token.Position{StartCol: 4, EndCol: 4}},
			{Type: token.IDENT, Literal: "int", Pos: token.Position{StartCol: 5, EndCol: 7}},
			{Type: token.LPAREN, Literal: "(", Pos: token.Position{StartCol: 8, EndCol: 8}},
			{Type: token.RPAREN, Literal: ")", Pos: token.Position{StartCol: 9, EndCol: 9}},
			{Type: token.RBRACES, Literal: "}}", Pos: token.Position{StartCol: 11, EndCol: 12}},
			{Type: token.EOF, Literal: "", Pos: token.Position{StartCol: 13, EndCol: 13}},
		})
	})
}

func TestForLoopStatement(t *testing.T) {
	inp := `@for(i = 0; i < 10; i++)`

	TokenizeString(t, inp, []token.Token{
		{Type: token.FOR, Literal: "@for", Pos: token.Position{EndCol: 3}},
		{Type: token.LPAREN, Literal: "(", Pos: token.Position{StartCol: 4, EndCol: 4}},
		{Type: token.IDENT, Literal: "i", Pos: token.Position{StartCol: 5, EndCol: 5}},
		{Type: token.ASSIGN, Literal: "=", Pos: token.Position{StartCol: 7, EndCol: 7}},
		{Type: token.INT, Literal: "0", Pos: token.Position{StartCol: 9, EndCol: 9}},
		{Type: token.SEMI, Literal: ";", Pos: token.Position{StartCol: 10, EndCol: 10}},
		{Type: token.IDENT, Literal: "i", Pos: token.Position{StartCol: 12, EndCol: 12}},
		{Type: token.LTHAN, Literal: "<", Pos: token.Position{StartCol: 14, EndCol: 14}},
		{Type: token.INT, Literal: "10", Pos: token.Position{StartCol: 16, EndCol: 17}},
		{Type: token.SEMI, Literal: ";", Pos: token.Position{StartCol: 18, EndCol: 18}},
		{Type: token.IDENT, Literal: "i", Pos: token.Position{StartCol: 20, EndCol: 20}},
		{Type: token.INC, Literal: "++", Pos: token.Position{StartCol: 21, EndCol: 22}},
		{Type: token.RPAREN, Literal: ")", Pos: token.Position{StartCol: 23, EndCol: 23}},
		{Type: token.EOF, Literal: "", Pos: token.Position{StartCol: 24, EndCol: 24}},
	})
}

func TestEachLoopStatement(t *testing.T) {
	inp := `@each(n in [1, 2, 3])@if(n == 2)@break@end{{ n }}@end`

	TokenizeString(t, inp, []token.Token{
		{Type: token.EACH, Literal: "@each", Pos: token.Position{EndCol: 4}},
		{Type: token.LPAREN, Literal: "(", Pos: token.Position{StartCol: 5, EndCol: 5}},
		{Type: token.IDENT, Literal: "n", Pos: token.Position{StartCol: 6, EndCol: 6}},
		{Type: token.IN, Literal: "in", Pos: token.Position{StartCol: 8, EndCol: 9}},
		{Type: token.LBRACKET, Literal: "[", Pos: token.Position{StartCol: 11, EndCol: 11}},
		{Type: token.INT, Literal: "1", Pos: token.Position{StartCol: 12, EndCol: 12}},
		{Type: token.COMMA, Literal: ",", Pos: token.Position{StartCol: 13, EndCol: 13}},
		{Type: token.INT, Literal: "2", Pos: token.Position{StartCol: 15, EndCol: 15}},
		{Type: token.COMMA, Literal: ",", Pos: token.Position{StartCol: 16, EndCol: 16}},
		{Type: token.INT, Literal: "3", Pos: token.Position{StartCol: 18, EndCol: 18}},
		{Type: token.RBRACKET, Literal: "]", Pos: token.Position{StartCol: 19, EndCol: 19}},
		{Type: token.RPAREN, Literal: ")", Pos: token.Position{StartCol: 20, EndCol: 20}},
		{Type: token.IF, Literal: "@if", Pos: token.Position{StartCol: 21, EndCol: 23}},
		{Type: token.LPAREN, Literal: "(", Pos: token.Position{StartCol: 24, EndCol: 24}},
		{Type: token.IDENT, Literal: "n", Pos: token.Position{StartCol: 25, EndCol: 25}},
		{Type: token.EQ, Literal: "==", Pos: token.Position{StartCol: 27, EndCol: 28}},
		{Type: token.INT, Literal: "2", Pos: token.Position{StartCol: 30, EndCol: 30}},
		{Type: token.RPAREN, Literal: ")", Pos: token.Position{StartCol: 31, EndCol: 31}},
		{Type: token.BREAK, Literal: "@break", Pos: token.Position{StartCol: 32, EndCol: 37}},
		{Type: token.END, Literal: "@end", Pos: token.Position{StartCol: 38, EndCol: 41}},
		{Type: token.LBRACES, Literal: "{{", Pos: token.Position{StartCol: 42, EndCol: 43}},
		{Type: token.IDENT, Literal: "n", Pos: token.Position{StartCol: 45, EndCol: 45}},
		{Type: token.RBRACES, Literal: "}}", Pos: token.Position{StartCol: 47, EndCol: 48}},
		{Type: token.END, Literal: "@end", Pos: token.Position{StartCol: 49, EndCol: 52}},
		{Type: token.EOF, Literal: "", Pos: token.Position{StartCol: 53, EndCol: 53}},
	})
}

func TestObjectStatement(t *testing.T) {
	inp := `{{ {"father": {"name": "John"}} }}`

	TokenizeString(t, inp, []token.Token{
		{Type: token.LBRACES, Literal: "{{", Pos: token.Position{EndCol: 1}},
		{Type: token.LBRACE, Literal: "{", Pos: token.Position{StartCol: 3, EndCol: 3}},
		{Type: token.STR, Literal: "father", Pos: token.Position{StartCol: 4, EndCol: 11}},
		{Type: token.COLON, Literal: ":", Pos: token.Position{StartCol: 12, EndCol: 12}},
		{Type: token.LBRACE, Literal: "{", Pos: token.Position{StartCol: 14, EndCol: 14}},
		{Type: token.STR, Literal: "name", Pos: token.Position{StartCol: 15, EndCol: 20}},
		{Type: token.COLON, Literal: ":", Pos: token.Position{StartCol: 21, EndCol: 21}},
		{Type: token.STR, Literal: "John", Pos: token.Position{StartCol: 23, EndCol: 28}},
		{Type: token.RBRACE, Literal: "}", Pos: token.Position{StartCol: 29, EndCol: 29}},
		{Type: token.RBRACE, Literal: "}", Pos: token.Position{StartCol: 30, EndCol: 30}},
		{Type: token.RBRACES, Literal: "}}", Pos: token.Position{StartCol: 32, EndCol: 33}},
		{Type: token.EOF, Literal: "", Pos: token.Position{StartCol: 34, EndCol: 34}},
	})
}

func TestBreakDirectives(t *testing.T) {
	inp := `@breakIf(true) @break @continue @continueIf(false)`

	TokenizeString(t, inp, []token.Token{
		{Type: token.BREAK_IF, Literal: "@breakIf", Pos: token.Position{EndCol: 7}},
		{Type: token.LPAREN, Literal: "(", Pos: token.Position{StartCol: 8, EndCol: 8}},
		{Type: token.TRUE, Literal: "true", Pos: token.Position{StartCol: 9, EndCol: 12}},
		{Type: token.RPAREN, Literal: ")", Pos: token.Position{StartCol: 13, EndCol: 13}},
		{Type: token.HTML, Literal: " ", Pos: token.Position{StartCol: 14, EndCol: 14}},
		{Type: token.BREAK, Literal: "@break", Pos: token.Position{StartCol: 15, EndCol: 20}},
		{Type: token.HTML, Literal: " ", Pos: token.Position{StartCol: 21, EndCol: 21}},
		{Type: token.CONTINUE, Literal: "@continue", Pos: token.Position{StartCol: 22, EndCol: 30}},
		{Type: token.HTML, Literal: " ", Pos: token.Position{StartCol: 31, EndCol: 31}},
		{Type: token.CONTINUE_IF, Literal: "@continueIf", Pos: token.Position{StartCol: 32, EndCol: 42}},
		{Type: token.LPAREN, Literal: "(", Pos: token.Position{StartCol: 43, EndCol: 43}},
		{Type: token.FALSE, Literal: "false", Pos: token.Position{StartCol: 44, EndCol: 48}},
		{Type: token.RPAREN, Literal: ")", Pos: token.Position{StartCol: 49, EndCol: 49}},
		{Type: token.EOF, Literal: "", Pos: token.Position{StartCol: 50, EndCol: 50}},
	})
}

func TestComponentDirective(t *testing.T) {
	inp := `@component("components/book-card", { c: card })`

	TokenizeString(t, inp, []token.Token{
		{Type: token.COMPONENT, Literal: "@component", Pos: token.Position{EndCol: 9}},
		{Type: token.LPAREN, Literal: "(", Pos: token.Position{StartCol: 10, EndCol: 10}},
		{Type: token.STR, Literal: "components/book-card", Pos: token.Position{StartCol: 11, EndCol: 32}},
		{Type: token.COMMA, Literal: ",", Pos: token.Position{StartCol: 33, EndCol: 33}},
		{Type: token.LBRACE, Literal: "{", Pos: token.Position{StartCol: 35, EndCol: 35}},
		{Type: token.IDENT, Literal: "c", Pos: token.Position{StartCol: 37, EndCol: 37}},
		{Type: token.COLON, Literal: ":", Pos: token.Position{StartCol: 38, EndCol: 38}},
		{Type: token.IDENT, Literal: "card", Pos: token.Position{StartCol: 40, EndCol: 43}},
		{Type: token.RBRACE, Literal: "}", Pos: token.Position{StartCol: 45, EndCol: 45}},
		{Type: token.RPAREN, Literal: ")", Pos: token.Position{StartCol: 46, EndCol: 46}},
		{Type: token.EOF, Literal: "", Pos: token.Position{StartCol: 47, EndCol: 47}},
	})
}

func TestComponentSlotDirective(t *testing.T) {
	t.Run("slots with parentheses", func(t *testing.T) {
		inp := `@component("card")@slot("top")<h1>Hello</h1>@end@end`

		TokenizeString(t, inp, []token.Token{
			{Type: token.COMPONENT, Literal: "@component", Pos: token.Position{EndCol: 9}},
			{Type: token.LPAREN, Literal: "(", Pos: token.Position{StartCol: 10, EndCol: 10}},
			{Type: token.STR, Literal: "card", Pos: token.Position{StartCol: 11, EndCol: 16}},
			{Type: token.RPAREN, Literal: ")", Pos: token.Position{StartCol: 17, EndCol: 17}},
			{Type: token.SLOT, Literal: "@slot", Pos: token.Position{StartCol: 18, EndCol: 22}},
			{Type: token.LPAREN, Literal: "(", Pos: token.Position{StartCol: 23, EndCol: 23}},
			{Type: token.STR, Literal: "top", Pos: token.Position{StartCol: 24, EndCol: 28}},
			{Type: token.RPAREN, Literal: ")", Pos: token.Position{StartCol: 29, EndCol: 29}},
			{Type: token.HTML, Literal: "<h1>Hello</h1>", Pos: token.Position{StartCol: 30, EndCol: 43}},
			{Type: token.END, Literal: "@end", Pos: token.Position{StartCol: 44, EndCol: 47}},
			{Type: token.END, Literal: "@end", Pos: token.Position{StartCol: 48, EndCol: 51}},
			{Type: token.EOF, Literal: "", Pos: token.Position{StartCol: 52, EndCol: 52}},
		})
	})

	t.Run("slots without parentheses", func(t *testing.T) {
		inp := `@component("card")@slotNICE@end@end`

		TokenizeString(t, inp, []token.Token{
			{Type: token.COMPONENT, Literal: "@component", Pos: token.Position{EndCol: 9}},
			{Type: token.LPAREN, Literal: "(", Pos: token.Position{StartCol: 10, EndCol: 10}},
			{Type: token.STR, Literal: "card", Pos: token.Position{StartCol: 11, EndCol: 16}},
			{Type: token.RPAREN, Literal: ")", Pos: token.Position{StartCol: 17, EndCol: 17}},
			{Type: token.SLOT, Literal: "@slot", Pos: token.Position{StartCol: 18, EndCol: 22}},
			{Type: token.HTML, Literal: "NICE", Pos: token.Position{StartCol: 23, EndCol: 26}},
			{Type: token.END, Literal: "@end", Pos: token.Position{StartCol: 27, EndCol: 30}},
			{Type: token.END, Literal: "@end", Pos: token.Position{StartCol: 31, EndCol: 34}},
			{Type: token.EOF, Literal: "", Pos: token.Position{StartCol: 35, EndCol: 35}},
		})
	})
}

// Comments should be ignored by the lexer
func TestCommentStatement(t *testing.T) {
	t.Run("Simple comment", func(t *testing.T) {
		inp := `<div>{{-- This is a comment --}}</div>`

		TokenizeString(t, inp, []token.Token{
			{Type: token.HTML, Literal: "<div>", Pos: token.Position{EndCol: 4}},
			{Type: token.HTML, Literal: "</div>", Pos: token.Position{StartCol: 32, EndCol: 37}},
			{Type: token.EOF, Literal: "", Pos: token.Position{StartCol: 38, EndCol: 38}},
		})
	})

	t.Run("Commented code", func(t *testing.T) {
		inp := `{{-- @each(u in users){{ u }}@end --}}`

		TokenizeString(t, inp, []token.Token{
			{Type: token.EOF, Literal: "", Pos: token.Position{StartCol: 38, EndCol: 38}},
		})
	})
}

func TestLexerCanReadIllegalDirectives(t *testing.T) {
	inp := `@if(false)@dump(@end`

	TokenizeString(t, inp, []token.Token{
		{Type: token.IF, Literal: "@if", Pos: token.Position{EndCol: 2}},
		{Type: token.LPAREN, Literal: "(", Pos: token.Position{StartCol: 3, EndCol: 3}},
		{Type: token.FALSE, Literal: "false", Pos: token.Position{StartCol: 4, EndCol: 8}},
		{Type: token.RPAREN, Literal: ")", Pos: token.Position{StartCol: 9, EndCol: 9}},
		{Type: token.DUMP, Literal: "@dump", Pos: token.Position{StartCol: 10, EndCol: 14}},
		{Type: token.LPAREN, Literal: "(", Pos: token.Position{StartCol: 15, EndCol: 15}},
		{Type: token.END, Literal: "@end", Pos: token.Position{StartCol: 16, EndCol: 19}},
		{Type: token.EOF, Literal: "", Pos: token.Position{StartCol: 20, EndCol: 20}},
	})
}
