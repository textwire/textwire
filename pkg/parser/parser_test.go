package parser

import (
	"reflect"
	"strings"
	"testing"

	"github.com/textwire/textwire/v4/pkg/ast"
	"github.com/textwire/textwire/v4/pkg/fail"
	"github.com/textwire/textwire/v4/pkg/lexer"
	"github.com/textwire/textwire/v4/pkg/position"
	"github.com/textwire/textwire/v4/pkg/token"
)

func TestParseIdentifier(t *testing.T) {
	inp := "{{ myName }}"

	identExpr, err := parseEmbedded[*ast.IdentExpr](inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	if err := testIdentExpr(identExpr, "myName"); err != nil {
		t.Fatal(err)
	}

	if err := testToken(identExpr, token.IDENT); err != nil {
		t.Fatal(err)
	}
}

func TestParseExpressionStatement(t *testing.T) {
	inp := "{{ 3 / 2 }}"

	infixExpr, err := parseEmbedded[*ast.InfixExpr](inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	if err := testToken(infixExpr, token.INT); err != nil {
		t.Fatal(err)
	}

	err = testTokPosition(infixExpr.Pos(), &position.Pos{
		StartCol: 3,
		EndCol:   7,
	})

	if err != nil {
		t.Fatal(err)
	}
}

func TestParseIntExpr(t *testing.T) {
	inp := "{{ 234 }}"

	intExpr, err := parseEmbedded[*ast.IntExpr](inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	if err := testToken(intExpr, token.INT); err != nil {
		t.Fatal(err)
	}

	if err := testIntExpr(intExpr, 234); err != nil {
		t.Fatal(err)
	}

	err = testTokPosition(intExpr.Pos(), &position.Pos{
		StartCol: 3,
		EndCol:   5,
	})

	if err != nil {
		t.Fatal(err)
	}
}

func TestParseFloatExpr(t *testing.T) {
	inp := "{{ 2.34149 }}"

	floatExpr, err := parseEmbedded[*ast.FloatExpr](inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	if err := testToken(floatExpr, token.FLOAT); err != nil {
		t.Fatal(err)
	}

	if err := testFloatExpr(floatExpr, 2.34149); err != nil {
		t.Fatal(err)
	}

	err = testTokPosition(floatExpr.Pos(), &position.Pos{
		StartCol: 3,
		EndCol:   9,
	})

	if err != nil {
		t.Fatal(err)
	}
}

func TestParseNilExpr(t *testing.T) {
	inp := "{{ nil }}"

	nilExpr, err := parseEmbedded[*ast.NilExpr](inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	if err := testToken(nilExpr, token.NIL); err != nil {
		t.Fatal(err)
	}

	if err := testNilExpr(nilExpr); err != nil {
		t.Fatal(err)
	}

	err = testTokPosition(nilExpr.Pos(), &position.Pos{
		StartCol: 3,
		EndCol:   5,
	})

	if err != nil {
		t.Fatal(err)
	}
}

func TestParseStrExpr(t *testing.T) {
	cases := []struct {
		inp      string
		expect   string
		startCol uint
		endCol   uint
	}{
		{`{{ "Hello World" }}`, "Hello World", 3, 15},
		{`{{ "Serhii \"Cho\"" }}`, `Serhii "Cho"`, 3, 18},
		{`{{ 'Hello World' }}`, "Hello World", 3, 15},
		{`{{ "" }}`, "", 3, 4},
	}

	for _, tc := range cases {
		strExpr, err := parseEmbedded[*ast.StrExpr](tc.inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}

		if err := testStrExpr(strExpr, tc.expect); err != nil {
			t.Fatal(err)
		}

		err = testTokPosition(strExpr.Pos(), &position.Pos{
			StartCol: tc.startCol,
			EndCol:   tc.endCol,
		})

		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestStrConcatenation(t *testing.T) {
	inp := `{{ 'Serhii' + " Anna" }}`

	infixExpr, err := parseEmbedded[*ast.InfixExpr](inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	if err := testToken(infixExpr, token.STR); err != nil {
		t.Fatal(err)
	}

	if err := testInfixExpr(infixExpr, "Serhii", "+", " Anna"); err != nil {
		t.Fatal(err)
	}
}

func TestParseInfixExpression(t *testing.T) {
	inp := "{{ 5 + 2 }}"

	infixExpr, err := parseEmbedded[*ast.InfixExpr](inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	if err := testInfixExpr(infixExpr, 5, "+", 2); err != nil {
		t.Fatal(err)
	}

	if err := testToken(infixExpr, token.INT); err != nil {
		t.Fatal(err)
	}

	err = testTokPosition(infixExpr.Pos(), &position.Pos{
		StartCol: 3,
		EndCol:   7,
	})

	if err != nil {
		t.Fatal(err)
	}
}

func TestParseGroupedExpression(t *testing.T) {
	inp := "{{ (5 + 3) * 2 }}"

	infixExpr, err := parseEmbedded[*ast.InfixExpr](inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	err = testTokPosition(infixExpr.Pos(), &position.Pos{
		StartCol: 3,
		EndCol:   13,
	})

	if err != nil {
		t.Fatal(err)
	}

	if err := testToken(infixExpr, token.LPAREN); err != nil {
		t.Fatal(err)
	}

	if err := testIntExpr(infixExpr.Right, 2); err != nil {
		t.Fatal(err)
	}

	if infixExpr.Op != "*" {
		t.Fatalf("infixExpr.Op is not %s, got %s", "*", infixExpr.Op)
	}

	leftInfixExpr, ok := infixExpr.Left.(*ast.InfixExpr)
	if !ok {
		t.Fatalf("infixExpr.Left is not an InfixExpr, got %T", infixExpr.Left)
	}

	if err := testInfixExpr(leftInfixExpr, 5, "+", 3); err != nil {
		t.Fatal(err)
	}
}

func TestParseInfixExp(t *testing.T) {
	cases := []struct {
		inp    string
		left   any
		op     string
		right  any
		endCol uint
		expTok token.TokenType
	}{
		{"{{ 5 + 8 }}", 5, "+", 8, 7, token.INT},
		{"{{ 10 - 2 }}", 10, "-", 2, 8, token.INT},
		{"{{ 2 * 2 }}", 2, "*", 2, 7, token.INT},
		{"{{ 44 / 4 }}", 44, "/", 4, 8, token.INT},
		{"{{ 5 % 4 }}", 5, "%", 4, 7, token.INT},
		{`{{ "me" + "her" }}`, "me", "+", "her", 14, token.STR},
		{`{{ 14 == 14 }}`, 14, "==", 14, 10, token.INT},
		{`{{ 10 != 1 }}`, 10, "!=", 1, 9, token.INT},
		{`{{ 19 > 31 }}`, 19, ">", 31, 9, token.INT},
		{`{{ 20 < 11 }}`, 20, "<", 11, 9, token.INT},
		{`{{ 19 >= 31 }}`, 19, ">=", 31, 10, token.INT},
		{`{{ 20 <= 11 }}`, 20, "<=", 11, 10, token.INT},
		{`{{ true && true }}`, true, "&&", true, 14, token.TRUE},
		{`{{ false || false }}`, false, "||", false, 16, token.FALSE},
	}

	for _, tc := range cases {
		infixExpr, err := parseEmbedded[*ast.InfixExpr](tc.inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}

		if err := testToken(infixExpr, tc.expTok); err != nil {
			t.Fatal(err)
		}

		err = testTokPosition(infixExpr.Pos(), &position.Pos{
			StartCol: 3,
			EndCol:   tc.endCol,
		})

		if err != nil {
			t.Fatal(err)
		}

		if err := testInfixExpr(infixExpr, tc.left, tc.op, tc.right); err != nil {
			t.Fatal(err)
		}
	}
}

func TestParseBooleanExpression(t *testing.T) {
	cases := []struct {
		inp      string
		expect   bool
		startCol uint
		endCol   uint
		expTok   token.TokenType
	}{
		{"{{ true }}", true, 3, 6, token.TRUE},
		{"{{ false }}", false, 3, 7, token.FALSE},
	}

	for _, tc := range cases {
		boolExpr, err := parseEmbedded[*ast.BoolExpr](tc.inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}

		if err := testToken(boolExpr, tc.expTok); err != nil {
			t.Fatal(err)
		}

		if err := testBoolExpr(boolExpr, tc.expect); err != nil {
			t.Fatal(err)
		}

		err = testTokPosition(boolExpr.Pos(), &position.Pos{
			StartCol: tc.startCol,
			EndCol:   tc.endCol,
		})

		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestParsePrefixExp(t *testing.T) {
	cases := []struct {
		inp    string
		op     string
		val    any
		endCol uint
		expTok token.TokenType
	}{
		{"{{ -5 }}", "-", 5, 4, token.SUB},
		{"{{ -10 }}", "-", 10, 5, token.SUB},
		{"{{ !true }}", "!", true, 7, token.NOT},
		{"{{ !false }}", "!", false, 8, token.NOT},
		{`{{ !"" }}`, "!", "", 5, token.NOT},
		{`{{ !0 }}`, "!", 0, 4, token.NOT},
		{`{{ -0 }}`, "-", 0, 4, token.SUB},
		{`{{ -0.0 }}`, "-", 0.0, 6, token.SUB},
		{`{{ !0.0 }}`, "!", 0.0, 6, token.NOT},
	}

	for _, tc := range cases {
		prefixExpr, err := parseEmbedded[*ast.PrefixExpr](tc.inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}

		if err := testToken(prefixExpr, tc.expTok); err != nil {
			t.Fatal(err)
		}

		err = testTokPosition(prefixExpr.Pos(), &position.Pos{
			StartCol: 3,
			EndCol:   tc.endCol,
		})

		if err != nil {
			t.Fatal(err)
		}

		if prefixExpr.Op != tc.op {
			t.Fatalf("prefixExpr.Op is not %s, got %s", tc.op, prefixExpr.Op)
		}

		if err := testLiteralExpr(prefixExpr.Right, tc.val); err != nil {
			t.Fatal(err)
		}
	}
}

func TestParseOperatorPrecedence(t *testing.T) {
	cases := []struct {
		id     uint
		inp    string
		expect string
	}{
		{
			id:     10,
			inp:    "{{ 1 * 2 }}",
			expect: "{{ (1 * 2) }}",
		},
		{
			id:     20,
			inp:    "<h2>{{ -2 + 3 }}</h2>",
			expect: "<h2>{{ ((-2) + 3) }}</h2>",
		},
		{
			id:     30,
			inp:    "{{ a + b + c }}",
			expect: "{{ ((a + b) + c) }}",
		},
		{
			id:     40,
			inp:    "{{ a + b / c }}",
			expect: "{{ (a + (b / c)) }}",
		},
		{
			id:     50,
			inp:    "{{ -2.float() }}",
			expect: "{{ (-(2.float())) }}",
		},
		{
			id:     60,
			inp:    "{{ -5.0.int() }}",
			expect: "{{ (-(5.0.int())) }}",
		},
		{
			id:     70,
			inp:    "{{ -obj.test }}",
			expect: "{{ (-(obj.test)) }}",
		},
		{
			id:     80,
			inp:    "{{ true && true || false }}",
			expect: "{{ ((true && true) || false) }}",
		},
		{
			id:     90,
			inp:    "{{ true ? 1 : 0 }}",
			expect: "{{ (true ? 1 : 0) }}",
		},
		{
			id:     100,
			inp:    "{{ true && false ? 1 : 0 }}",
			expect: "{{ ((true && false) ? 1 : 0) }}",
		},
		{
			id:     110,
			inp:    "{{ true && false || 1 ? 1 : 0 }}",
			expect: "{{ (((true && false) || 1) ? 1 : 0) }}",
		},
		{
			id:     120,
			inp:    "{{ -2.float() && -2.0.int() ? 1 : 0 }}",
			expect: "{{ (((-(2.float())) && (-(2.0.int()))) ? 1 : 0) }}",
		},
		{
			id:     130,
			inp:    "{{ !defined(age) || !defined(name) ? 1 : 0 }}",
			expect: "{{ (((!(defined(age))) || (!(defined(name)))) ? 1 : 0) }}",
		},
		{
			id:     140,
			inp:    "{{ defined(name) }}",
			expect: "{{ (defined(name)) }}",
		},
		{
			id:     150,
			inp:    "{{ long = user.name.len() > 0 }}",
			expect: "{{ (long = (((user.name).len()) > 0)) }}",
		},
		{
			id:     160,
			inp:    "{{ user && user.name == 'serhii' }}",
			expect: `{{ (user && ((user.name) == "serhii")) }}`,
		},
	}

	for _, tc := range cases {
		l := lexer.New(tc.inp)
		p := New(l, nil)
		prog := p.ParseProgram()

		if p.HasErrors() {
			t.Fatalf("Case: %d. Parser error:  %v", tc.id, p.Errors()[0])
		}

		if prog.String() != tc.expect {
			t.Fatalf("Case: %d. Expect %s but got %s", tc.id, tc.expect, prog)
		}
	}
}

func TestErrorHandling(t *testing.T) {
	cases := []struct {
		id  uint
		inp string
		err *fail.Error
	}{
		// Embedded chunk
		{
			id:  10,
			inp: `{{ obj."str" }}`,
			err: fail.New(
				&position.Pos{StartCol: 7, EndCol: 11},
				"",
				fail.OriginPars,
				fail.ErrWrongPeekToken,
				token.String(token.IDENT),
				"str",
			),
		},
		{
			id:  11,
			inp: `{{ { "1st": "nice" }.1st }}`,
			err: fail.New(
				&position.Pos{StartCol: 21, EndCol: 21},
				"",
				fail.OriginPars,
				fail.ErrWrongPeekToken,
				token.String(token.IDENT),
				"1",
			),
		},
		{
			id:  12,
			inp: "{{ true ? 100 }}",
			err: fail.New(
				&position.Pos{StartCol: 14, EndCol: 15},
				"",
				fail.OriginPars,
				fail.ErrWrongPeekToken,
				token.String(token.COLON),
				"}}",
			),
		},
		{
			id:  20,
			inp: "{{ 5 + }}",
			err: fail.New(
				&position.Pos{StartCol: 3, EndCol: 5},
				"",
				fail.OriginPars,
				fail.ErrExpectExprAfter,
				"+",
			),
		},
		{
			id:  30,
			inp: "{{ myVar = }}",
			err: fail.New(
				&position.Pos{StartCol: 11, EndCol: 12},
				"",
				fail.OriginPars,
				fail.ErrExpectExprAfter,
				"=",
			),
		},
		{
			id:  40,
			inp: "{{ }}",
			err: fail.New(&position.Pos{EndCol: 4}, "", fail.OriginPars, fail.ErrEmptyBraces),
		},
		{
			id:  50,
			inp: `{{ 1 ~ 8 }}`,
			err: fail.New(
				&position.Pos{StartCol: 5, EndCol: 5},
				"",
				fail.OriginPars,
				fail.ErrIllegalToken,
				"~",
			),
		},
		{
			id:  60,
			inp: "{{ ) }}",
			err: fail.New(
				&position.Pos{StartCol: 3, EndCol: 3},
				"",
				fail.OriginPars,
				fail.ErrIllegalToken,
				token.String(token.RPAREN),
			),
		},
		// Use directive
		{
			id:  70,
			inp: "@use('')",
			err: fail.New(
				&position.Pos{StartCol: 5, EndCol: 6},
				"",
				fail.OriginPars,
				fail.ErrNameCannotBeEmpty,
				"@use",
			),
		},
		{
			id:  71,
			inp: "@use('base') @use('main')",
			err: fail.New(
				&position.Pos{StartCol: 13, EndCol: 24},
				"",
				fail.OriginPars,
				fail.ErrOnlyOneUseDir,
			),
		},
		{
			id:  80,
			inp: "@use(1)",
			err: fail.New(
				&position.Pos{StartCol: 5, EndCol: 5},
				"",
				fail.OriginPars,
				fail.ErrWrongTokenType,
				token.String(token.STR),
				token.String(token.INT),
			),
		},
		{
			id:  90,
			inp: `@use "name"`,
			err: fail.New(
				&position.Pos{StartCol: 4, EndCol: 10},
				"",
				fail.OriginPars,
				fail.ErrWrongPeekToken,
				token.String(token.LPAREN),
				` "name"`,
			),
		},
		// Component
		{
			id:  200,
			inp: "@component('')@end",
			err: fail.New(
				&position.Pos{StartCol: 11, EndCol: 12},
				"",
				fail.OriginPars,
				fail.ErrNameCannotBeEmpty,
				"@component",
			),
		},
		{
			id:  210,
			inp: "@component(3.3)@end",
			err: fail.New(
				&position.Pos{StartCol: 11, EndCol: 13},
				"",
				fail.OriginPars,
				fail.ErrWrongTokenType,
				token.String(token.STR),
				token.String(token.FLOAT),
			),
		},
		{
			id:  220,
			inp: "@component('~user'",
			err: fail.New(
				&position.Pos{StartCol: 18, EndCol: 18},
				"",
				fail.OriginPars,
				fail.ErrWrongPeekToken,
				token.String(token.RPAREN),
				"",
			),
		},
		{
			id:  230,
			inp: "@component   ('",
			err: fail.New(
				&position.Pos{StartCol: 14, EndCol: 15},
				"",
				fail.OriginPars,
				fail.ErrNameCannotBeEmpty,
				"@component",
			),
		},
		{
			id:  240,
			inp: "@component",
			err: fail.New(
				&position.Pos{StartCol: 10, EndCol: 10},
				"",
				fail.OriginPars,
				fail.ErrWrongPeekToken,
				token.String(token.LPAREN),
				"",
			),
		},
		// Reserve
		{
			id:  300,
			inp: "@reserve('title')\n@reserve('title')",
			err: fail.New(
				&position.Pos{StartLine: 1, StartCol: 9, EndLine: 1, EndCol: 15},
				"",
				fail.OriginPars,
				fail.ErrDuplicateReserves,
				"title",
				"",
			),
		},
		{
			id:  310,
			inp: "@reserve(1)",
			err: fail.New(
				&position.Pos{StartCol: 9, EndCol: 9},
				"",
				fail.OriginPars,
				fail.ErrWrongTokenType,
				token.String(token.STR),
				token.String(token.INT),
			),
		},
		{
			id:  320,
			inp: "@reserve('')",
			err: fail.New(
				&position.Pos{StartCol: 9, EndCol: 10},
				"",
				fail.OriginPars,
				fail.ErrNameCannotBeEmpty,
				"@reserve",
			),
		},
		// Insert
		{
			id:  400,
			inp: "@insert('x', 'x')\n@insert('y', 'y')\n@insert('x', 'x')",
			err: fail.New(
				&position.Pos{StartLine: 2, StartCol: 8, EndLine: 2, EndCol: 10},
				"",
				fail.OriginPars,
				fail.ErrDuplicateInserts,
				"x",
			),
		},
		{
			id:  410,
			inp: "@insert('', 'x')",
			err: fail.New(
				&position.Pos{StartCol: 8, EndCol: 9},
				"",
				fail.OriginPars,
				fail.ErrNameCannotBeEmpty,
				"@insert",
			),
		},
		{
			id:  420,
			inp: "@insert([1, 2], 'test')",
			err: fail.New(
				&position.Pos{StartCol: 8, EndCol: 8},
				"",
				fail.OriginPars,
				fail.ErrWrongTokenType,
				token.String(token.STR),
				token.String(token.LBRACKET),
			),
		},
		{
			id:  430,
			inp: "@insert('nice",
			err: fail.New(
				&position.Pos{StartCol: 14, EndCol: 14},
				"",
				fail.OriginPars,
				fail.ErrWrongPeekToken,
				token.String(token.RPAREN),
				"",
			),
		},
		{
			id:  440,
			inp: "@insert ('nice'",
			err: fail.New(
				&position.Pos{StartCol: 15, EndCol: 15},
				"",
				fail.OriginPars,
				fail.ErrWrongPeekToken,
				token.String(token.RPAREN),
				"",
			),
		},
		{
			id:  450,
			inp: "@insert('nice'@end",
			err: fail.New(
				&position.Pos{StartCol: 14, EndCol: 17},
				"",
				fail.OriginPars,
				fail.ErrWrongPeekToken,
				token.String(token.RPAREN),
				"@end",
			),
		},
		{
			id:  460,
			inp: "@insert    ('nice' {{ 'nice' }}@end",
			err: fail.New(
				&position.Pos{StartCol: 19, EndCol: 20},
				"",
				fail.OriginPars,
				fail.ErrWrongPeekToken,
				token.String(token.RPAREN),
				"{{",
			),
		},
		{
			id:  470,
			inp: "@insert('title')<h3>Hello</h3>",
			err: fail.New(
				&position.Pos{EndCol: 6},
				"",
				fail.OriginPars,
				fail.ErrInsertMustHaveContent,
				"title",
			),
		},
		// For directive
		{
			id:  500,
			inp: "@for(i = 0; i < 4; i+2){{ i }}@end",
			err: fail.New(
				&position.Pos{StartCol: 19, EndCol: 21},
				"",
				fail.OriginPars,
				fail.ErrForLoopExpectStmt,
				"(i + 2)",
			),
		},
		{
			id:  600,
			inp: "@for(c = 0.0; c < 4.0; c+2.0){{ c }}@end",
			err: fail.New(
				&position.Pos{StartCol: 23, EndCol: 27},
				"",
				fail.OriginPars,
				fail.ErrForLoopExpectStmt,
				"(c + 2.0)",
			),
		},
		{
			id:  610,
			inp: "@for(;;    true  ){{ c }}@end",
			err: fail.New(
				&position.Pos{StartCol: 11, EndCol: 14},
				"",
				fail.OriginPars,
				fail.ErrForLoopExpectStmt,
				"true",
			),
		},
		// Each directive
		{
			id:  700,
			inp: "@each( {{ 'nice' }}@end",
			err: fail.New(
				&position.Pos{StartCol: 7, EndCol: 8},
				"",
				fail.OriginPars,
				fail.ErrWrongTokenType,
				token.String(token.IDENT),
				token.String(token.LBRACES),
			),
		},
		// If directive
		{
			id:  800,
			inp: "@if( {{ 'nice' }}@end",
			err: fail.New(
				&position.Pos{StartCol: 5, EndCol: 6},
				"",
				fail.OriginPars,
				fail.ErrIllegalToken,
				token.String(token.LBRACES),
			),
		},
		{
			id:  810,
			inp: "@if(false)@dump(@end",
			err: fail.New(
				&position.Pos{StartCol: 16, EndCol: 19},
				"",
				fail.OriginPars,
				fail.ErrExpectExprAfter,
				token.String(token.LPAREN),
			),
		},
		// Global functions
		{
			id:  900,
			inp: "{{ defined() }}",
			err: fail.New(
				&position.Pos{StartCol: 3, EndCol: 11},
				"",
				fail.OriginPars,
				fail.ErrGlobalFuncFewArgs,
				"defined",
				1,
				0,
			),
		},
		{
			id:  910,
			inp: "{{ hasValue() }}",
			err: fail.New(
				&position.Pos{StartCol: 3, EndCol: 12},
				"",
				fail.OriginPars,
				fail.ErrGlobalFuncFewArgs,
				"hasValue",
				1,
				0,
			),
		},
		{
			id:  920,
			inp: "{{ formatDate() }}",
			err: fail.New(
				&position.Pos{StartCol: 3, EndCol: 14},
				"",
				fail.OriginPars,
				fail.ErrGlobalFuncFewArgs,
				"formatDate",
				2,
				0,
			),
		},
		{
			id:  930,
			inp: "{{ formatDate(str) }}",
			err: fail.New(
				&position.Pos{StartCol: 3, EndCol: 17},
				"",
				fail.OriginPars,
				fail.ErrGlobalFuncFewArgs,
				"formatDate",
				2,
				1,
			),
		},
		{
			id:  940,
			inp: "{{ formatDate(str, str2, str3) }}",
			err: fail.New(
				&position.Pos{StartCol: 3, EndCol: 29},
				"",
				fail.OriginPars,
				fail.ErrGlobalFuncLotsOfArgs,
				"formatDate",
				2,
				3,
			),
		},
		{
			id:  950,
			inp: "@component('name')@pass('')<content>@end@end",
			err: fail.New(
				&position.Pos{StartCol: 24, EndCol: 25},
				"",
				fail.OriginPars,
				fail.ErrNameCannotBeEmpty,
				"@pass",
			),
		},
		{
			id:  960,
			inp: "<div>@slot('name')Nice@slot('name')</div>",
			err: fail.New(
				&position.Pos{StartCol: 28, EndCol: 33},
				"",
				fail.OriginPars,
				fail.ErrDuplicateSlots,
				"name",
				"",
			),
		},
	}

	for _, tc := range cases {
		l := lexer.New(tc.inp)
		p := New(l, nil)
		p.ParseProgram()

		if !p.HasErrors() {
			t.Fatalf("Case: %d. No errors found in input %q", tc.id, tc.inp)
		}

		err := p.Errors()[0]
		if err.String() != tc.err.String() {
			t.Fatalf("Case: %d. Expect error message:\n%q\ngot:\n%q", tc.id, tc.err, err)
		}

		if !reflect.DeepEqual(err.Pos(), tc.err.Pos()) {
			t.Fatalf(
				"Case: %d. Wrong position on error message, expect %v, got: %v",
				tc.id,
				tc.err.Pos(),
				err.Pos(),
			)
		}
	}
}

func TestParseTernaryExpr(t *testing.T) {
	inp := `{{ true ? 100 : "Some string" }}`

	terExpr, err := parseEmbedded[*ast.TernaryExpr](inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	if err := testToken(terExpr, token.TRUE); err != nil {
		t.Fatal(err)
	}

	err = testTokPosition(terExpr.Pos(), &position.Pos{
		StartCol: 3,
		EndCol:   28,
	})

	if err != nil {
		t.Fatal(err)
	}

	if err := testBoolExpr(terExpr.Cond, true); err != nil {
		t.Fatal(err)
	}

	if err := testIntExpr(terExpr.IfExpr, 100); err != nil {
		t.Fatal(err)
	}

	if err := testStrExpr(terExpr.ElseExpr, "Some string"); err != nil {
		t.Fatal(err)
	}
}

func TestParseIfDir(t *testing.T) {
	t.Run("regular @if", func(t *testing.T) {
		inp := `@if(true)1@end`

		ifDir, err := parseDirective[*ast.IfDir](inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}

		if err := testToken(ifDir, token.IF); err != nil {
			t.Fatal(err)
		}

		err = testTokPosition(ifDir.Pos(), &position.Pos{
			StartCol: 0,
			EndCol:   13,
		})

		if err != nil {
			t.Fatal(err)
		}

		if err := testIfDir(ifDir, true, "1"); err != nil {
			t.Fatal(err)
		}

		if ifDir.ElseBlock != nil {
			t.Fatalf("ifDir.ElseBlock is not nil, got %T", ifDir.ElseBlock)
		}

		if len(ifDir.ElseifDirs) != 0 {
			t.Fatalf("ifDir.ElseifDirs is not empty, got %d", len(ifDir.ElseifDirs))
		}
	})

	t.Run("@if with @else", func(t *testing.T) {
		inp := `@if(true)1@else2@end`

		ifDir, err := parseDirective[*ast.IfDir](inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}

		if err := testToken(ifDir, token.IF); err != nil {
			t.Fatal(err)
		}

		if err := testToken(ifDir.IfBlock, token.TEXT); err != nil {
			t.Fatal(err)
		}

		if err := testToken(ifDir.ElseBlock, token.TEXT); err != nil {
			t.Fatal(err)
		}

		if err := testIfDir(ifDir, true, "1"); err != nil {
			t.Fatal(err)
		}

		if err := testBlock(ifDir.ElseBlock, "2"); err != nil {
			t.Fatal(err)
		}

		err = testTokPosition(ifDir.Pos(), &position.Pos{
			StartCol: 0,
			EndCol:   19,
		})

		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("nested @if with @else", func(t *testing.T) {
		inp := `@if(true)
			@if(false)
				James
			@elseif(false)
				John
			@else
				@if(true){{ "Marry" }}@end
			@end
		@else
			@if(true)Anna@end
		@end`

		ifDir, err := parseDirective[*ast.IfDir](inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}

		if len(ifDir.IfBlock.Chunks) != 1 {
			t.Fatalf(
				"ifDir.IfBlock.Chunks does not contain 1 chunk1, got %d",
				len(ifDir.IfBlock.Chunks),
			)
		}

		if err := testToken(ifDir, token.IF); err != nil {
			t.Fatal(err)
		}

		if err := testToken(ifDir.IfBlock, token.TEXT); err != nil {
			t.Fatal(err)
		}

		if err := testToken(ifDir.ElseBlock, token.TEXT); err != nil {
			t.Fatal(err)
		}
	})
}

func TestParseIfElseIfDir(t *testing.T) {
	inp := `@if(true)first@elseif(false)second@end`

	ifDir, err := parseDirective[*ast.IfDir](inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	if err := testIfDir(ifDir, true, "first"); err != nil {
		t.Fatal(err)
	}

	if ifDir.ElseBlock != nil {
		t.Fatalf("ifDir.ElseBlock is not nil, got %T", ifDir.ElseBlock)
	}

	if len(ifDir.ElseifDirs) != 1 {
		t.Fatalf("ifDir.ElseifDirs does not contain 1 dir, got %d", len(ifDir.ElseifDirs))
	}

	elseifDir := ifDir.ElseifDirs[0]
	if err := testBoolExpr(elseifDir.Cond, false); err != nil {
		t.Fatal(err)
	}

	if len(elseifDir.Block.Chunks) != 1 {
		t.Fatalf(
			"elseifDir.Block.Chunks does not contain 1 chunk, got %d",
			len(elseifDir.Block.Chunks),
		)
	}

	text, ok := elseifDir.Block.Chunks[0].(*ast.Text)
	if !ok {
		t.Fatalf(
			"elseifDir.Block.Chunks[0] is not an Text, got %T",
			elseifDir.Block.Chunks[0],
		)
	}

	if text.String() != "second" {
		t.Fatalf("text.String() is not %q, got %q", "second", text)
	}
}

func TestParseElseIfWithElseDir(t *testing.T) {
	inp := `@if(true)1@elseif(false)2@else3@end`

	ifDir, err := parseDirective[*ast.IfDir](inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	if err := testIfDir(ifDir, true, "1"); err != nil {
		t.Fatal(err)
	}

	if err := testBlock(ifDir.ElseBlock, "3"); err != nil {
		t.Fatal(err)
	}

	if len(ifDir.ElseifDirs) != 1 {
		t.Fatalf(
			"ifDir.ElseifDirs does not contain 1 dir, got %d",
			len(ifDir.ElseifDirs),
		)
	}

	elseifDir := ifDir.ElseifDirs[0]
	if err := testBoolExpr(elseifDir.Cond, false); err != nil {
		t.Fatal(err)
	}

	if len(elseifDir.Block.Chunks) != 1 {
		t.Fatalf(
			"elseifDir.Block.Chunks does not contain 1 chunk, got %d",
			len(elseifDir.Block.Chunks),
		)
	}

	text, ok := elseifDir.Block.Chunks[0].(*ast.Text)
	if !ok {
		t.Fatalf(
			"elseifDir.Block.Chunks[0] is not an Text, got %T",
			elseifDir.Block.Chunks[0],
		)
	}

	if text.String() != "2" {
		t.Fatalf("text.String() is not %s, got %s", "2", text)
	}
}

func TestParseAssignStmt(t *testing.T) {
	cases := []struct {
		id       uint
		inp      string
		str      string
		startCol uint
		endCol   uint
	}{
		{
			id:       10,
			inp:      `{{ name = "Anna" }}`,
			str:      `(name = "Anna")`,
			startCol: 3,
			endCol:   15,
		},
		{
			id:       20,
			inp:      `{{ myAge = 34 }}`,
			str:      `(myAge = 34)`,
			startCol: 3,
			endCol:   12,
		},
		{
			id:       30,
			inp:      `{{ me.age = 34 }}`,
			str:      `((me.age) = 34)`,
			startCol: 3,
			endCol:   13,
		},
		{
			id:       40,
			inp:      `{{ arr[0] = 1 }}`,
			str:      `((arr[0]) = 1)`,
			startCol: 3,
			endCol:   12,
		},
		{
			id:       50,
			inp:      `{{ arr[234][2][23].name.first = "Anna" }}`,
			str:      `((((((arr[234])[2])[23]).name).first) = "Anna")`,
			startCol: 3,
			endCol:   37,
		},
		{
			id:       60,
			inp:      `{{ (obj.one.two) = "test" }}`,
			str:      `(((obj.one).two) = "test")`,
			startCol: 3,
			endCol:   24,
		},
	}

	for _, tc := range cases {
		assignStmt, err := parseEmbedded[*ast.AssignStmt](tc.inp, defaultParseOpts)
		if err != nil {
			t.Fatalf("Case: %d. %v", tc.id, err)
		}

		str := assignStmt.String()
		if str != tc.str {
			t.Fatalf("Case: %d. assignStmt.String() is not %s, got %s", tc.id, tc.str, str)
		}

		err = testTokPosition(assignStmt.Pos(), &position.Pos{
			StartCol: tc.startCol,
			EndCol:   tc.endCol,
		})

		if err != nil {
			t.Fatalf("Case: %d. %v", tc.id, err)
		}
	}
}

func TestParseUseDir(t *testing.T) {
	inp := `@use("main")`

	useDir, err := parseDirective[*ast.UseDir](inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	if useDir.Name.Val != "main" {
		t.Fatalf("useDir.Name.Val is not 'main', got %s", useDir.Name.Val)
	}

	if err := testToken(useDir, token.USE); err != nil {
		t.Fatal(err)
	}

	err = testTokPosition(useDir.Pos(), &position.Pos{
		StartCol: 0,
		EndCol:   11,
	})

	if err != nil {
		t.Fatal(err)
	}

	if useDir.LayoutProg != nil {
		t.Fatalf("useDir.LayoutProg is not nil, got %T", useDir.LayoutProg)
	}

	if useDir.String() != inp {
		t.Fatalf("useDir.String() is not %s, got %s", inp, useDir)
	}
}

func TestParseReserveDir(t *testing.T) {
	inp := `<div>@reserve("content")</div>`

	chunks, err := parseChunks(inp, parseOpts{chunksCount: 3, checkErrors: true})
	if err != nil {
		t.Fatal(err)
	}

	reserveDir, ok := chunks[1].(*ast.ReserveDir)
	if !ok {
		t.Fatalf("chunks[1] is not a ReserveDir, got %T", chunks[1])
	}

	if err := testToken(reserveDir, token.RESERVE); err != nil {
		t.Fatal(err)
	}

	err = testTokPosition(reserveDir.Pos(), &position.Pos{
		StartCol: 5,
		EndCol:   23,
	})

	if err != nil {
		t.Fatal(err)
	}

	if reserveDir.Name.Val != "content" {
		t.Fatalf("reserveDir.Name.Val is not 'content', got %s", reserveDir.Name.Val)
	}

	if reserveDir.String() == inp {
		t.Fatalf("reserveDir.String() is not %s, got %s", inp, reserveDir)
	}
}

func TestInsertDir(t *testing.T) {
	t.Run("@insert with block", func(t *testing.T) {
		inp := `<h1>@insert("content")<h1>Some content</h1>@end</h1>`

		chunks, err := parseChunks(inp, parseOpts{chunksCount: 3, checkErrors: true})
		if err != nil {
			t.Fatal(err)
		}

		insertDir, ok := chunks[1].(*ast.InsertDir)
		if !ok {
			t.Fatalf("chunks[1] is not an InsertDir, got %T", chunks[1])
		}

		if err := testToken(insertDir, token.INSERT); err != nil {
			t.Fatal(err)
		}

		err = testTokPosition(insertDir.Pos(), &position.Pos{
			StartCol: 4,
			EndCol:   46,
		})

		if err != nil {
			t.Fatal(err)
		}

		if insertDir.Name.Val != "content" {
			t.Fatalf("insertDir.Name.Val is not 'content', got %s", insertDir.Name.Val)
		}

		if insertDir.Block.String() != "<h1>Some content</h1>" {
			t.Fatalf(
				"insertDir.Block.String() is not '<h1>Some content</h1>', got %s",
				insertDir.Block,
			)
		}
	})

	t.Run("@insert with argument", func(t *testing.T) {
		inp := "<h1>@insert('content', 'Some content')</h1>"

		chunks, err := parseChunks(inp, parseOpts{chunksCount: 3, checkErrors: true})
		if err != nil {
			t.Fatal(err)
		}

		insertDir, ok := chunks[1].(*ast.InsertDir)
		if !ok {
			t.Fatalf("chunks[1] is not an InsertDir, got %T", chunks[1])
		}

		err = testTokPosition(insertDir.Pos(), &position.Pos{
			StartCol: 4,
			EndCol:   37,
		})

		if err != nil {
			t.Fatal(err)
		}

		if err := testToken(insertDir, token.INSERT); err != nil {
			t.Fatal(err)
		}

		if insertDir.Name.Val != "content" {
			t.Fatalf("insertDir.Name.Val is not 'content', got %s", insertDir.Name.Val)
		}

		if insertDir.Block != nil {
			t.Fatalf("insertDir.Block is not nil, got %T", insertDir.Block)
		}

		if insertDir.Argument.String() != `"Some content"` {
			t.Fatalf(
				"insertDir.Argument.String() is not 'Some content', got %s",
				insertDir.Argument,
			)
		}
	})
}

func TestParseArr(t *testing.T) {
	inp := `{{ [11, 234,] }}`

	arr, err := parseEmbedded[*ast.ArrExpr](inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	if err := testToken(arr, token.LBRACKET); err != nil {
		t.Fatal(err)
	}

	err = testTokPosition(arr.Pos(), &position.Pos{
		StartCol: 3,
		EndCol:   12,
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(arr.Elements) != 2 {
		t.Fatalf("len(arr.Elements) is not 2, got %d", len(arr.Elements))
	}

	if err := testIntExpr(arr.Elements[0], 11); err != nil {
		t.Fatal(err)
	}

	if err := testIntExpr(arr.Elements[1], 234); err != nil {
		t.Fatal(err)
	}
}

func TestParseIndexExp(t *testing.T) {
	inp := `{{ arr[1 + 2][2] }}`

	exp, err := parseEmbedded[*ast.IndexExpr](inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	if err := testToken(exp, token.IDENT); err != nil {
		t.Fatal(err)
	}

	err = testTokPosition(exp.Pos(), &position.Pos{
		StartCol: 3,
		EndCol:   15,
	})

	if err != nil {
		t.Fatal(err)
	}

	if exp.String() != "((arr[(1 + 2)])[2])" {
		t.Fatalf("indexExp.String() is not '(arr[(1 + 2)])', got %s", exp)
	}
}

func TestParseIncStmt(t *testing.T) {
	inp := "{{ i++ }}"

	stmt, err := parseEmbedded[*ast.IncStmt](inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	if err := testToken(stmt, token.INC); err != nil {
		t.Fatal(err)
	}

	if err := testIdentExpr(stmt.Left, "i"); err != nil {
		t.Fatal(err)
	}

	if stmt.String() != "(i++)" {
		t.Fatalf("stmt.String() is not '(i++)', got %s", stmt)
	}
}

func TestParseDecStmt(t *testing.T) {
	inp := "{{ i-- }}"

	stmt, err := parseEmbedded[*ast.DecStmt](inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	if err := testToken(stmt, token.DEC); err != nil {
		t.Fatal(err)
	}

	if err := testIdentExpr(stmt.Left, "i"); err != nil {
		t.Fatal(err)
	}

	if stmt.String() != "(i--)" {
		t.Fatalf("stmt.String() is not '(i--)', got %s", stmt)
	}
}

func TestParseTwoStatements(t *testing.T) {
	inp := `{{ name = "Anna"; name }}`

	segments, err := parseEmbeddedSegments(inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	if len(segments) != 2 {
		t.Fatalf("expected 2 segments, got %d", len(segments))
	}

	assignStmt, ok := segments[0].(*ast.AssignStmt)
	if !ok {
		t.Fatalf("segments[0] is not AssignStmt, got %T", segments[0])
	}

	if err := testIdentExpr(assignStmt.Left, "name"); err != nil {
		t.Fatal(err)
	}

	if err := testToken(assignStmt, token.IDENT); err != nil {
		t.Fatal(err)
	}

	nameExpr, ok := segments[1].(*ast.IdentExpr)
	if !ok {
		t.Fatalf("segments[1] is not IdentExpr, got %T", segments[1])
	}

	if err := testIdentExpr(nameExpr, "name"); err != nil {
		t.Fatal(err)
	}

	if err := testStrExpr(assignStmt.Right, "Anna"); err != nil {
		t.Fatal(err)
	}

	if assignStmt.String() != `(name = "Anna")` {
		t.Fatalf("assignStmt.String() is not 'name = \"Anna\"', got %s", assignStmt)
	}

	if nameExpr.String() != `name` {
		t.Fatalf("nameExpr.String() is not 'name', got %s", nameExpr)
	}
}

func TestParseTwoExpression(t *testing.T) {
	inp := `{{ 1; 2 }}`
	segments, err := parseEmbeddedSegments(inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	if len(segments) != 2 {
		t.Fatalf("expected 2 segments, got %d", len(segments))
	}

	exp1, ok := segments[0].(*ast.IntExpr)
	if !ok {
		t.Fatalf("segments[0] is not IntExpr, got %T", segments[0])
	}

	if err := testIntExpr(exp1, 1); err != nil {
		t.Fatal(err)
	}

	exp2, ok := segments[1].(*ast.IntExpr)
	if !ok {
		t.Fatalf("segments[1] is not IntExpr, got %T", segments[1])
	}

	if err := testIntExpr(exp2, 2); err != nil {
		t.Fatal(err)
	}
}

func TestParseGlobalCallExp(t *testing.T) {
	inp := `{{ defined(var1, var2) }}`

	globalCallExpr, err := parseEmbedded[*ast.GlobalCallExpr](inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	if err := testToken(globalCallExpr, token.IDENT); err != nil {
		t.Fatal(err)
	}

	if globalCallExpr.Name != ast.GlobalFuncName("defined") {
		t.Fatalf("Global function name must be 'defined', bot '%s'", globalCallExpr.Name)
	}

	if err := testIdentExpr(globalCallExpr.Arguments[0], "var1"); err != nil {
		t.Fatal(err)
	}

	if err := testIdentExpr(globalCallExpr.Arguments[1], "var2"); err != nil {
		t.Fatal(err)
	}

	err = testTokPosition(globalCallExpr.Pos(), &position.Pos{
		StartCol: 3,
		EndCol:   21,
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(globalCallExpr.Arguments) != 2 {
		t.Fatalf("len(globalCallExpr.Arguments) is not 2, got %d", len(globalCallExpr.Arguments))
	}
}

func TestParseCallExp(t *testing.T) {
	inp := `{{ "Serhii Cho".split(" ") }}`

	callExpr, err := parseEmbedded[*ast.CallExpr](inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	if err := testToken(callExpr, token.IDENT); err != nil {
		t.Fatal(err)
	}

	err = testTokPosition(callExpr.Pos(), &position.Pos{
		StartCol: 16,
		EndCol:   25,
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := testStrExpr(callExpr.Receiver, "Serhii Cho"); err != nil {
		t.Fatal(err)
	}

	if err := testIdentExpr(callExpr.Function, "split"); err != nil {
		t.Fatal(err)
	}

	if len(callExpr.Arguments) != 1 {
		t.Fatalf("len(callExpr.Arguments) is not 1, got %d", len(callExpr.Arguments))
	}

	if err := testStrExpr(callExpr.Arguments[0], " "); err != nil {
		t.Fatal(err)
	}
}

func TestParseCallExpWithExpressionList(t *testing.T) {
	inp := `{{ "nice".replace("n", "") }}`

	callExpr, err := parseEmbedded[*ast.CallExpr](inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	if err := testToken(callExpr, token.IDENT); err != nil {
		t.Fatal(err)
	}

	err = testTokPosition(callExpr.Pos(), &position.Pos{
		StartCol: 10,
		EndCol:   25,
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(callExpr.Arguments) != 2 {
		t.Fatalf("len(callExpr.Arguments) is not 2, got %d", len(callExpr.Arguments))
	}
}

func TestParseCallExpWithEmptyString(t *testing.T) {
	inp := `{{ "".len() }}`

	callExpr, err := parseEmbedded[*ast.CallExpr](inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	if err := testToken(callExpr, token.IDENT); err != nil {
		t.Fatal(err)
	}

	if err := testStrExpr(callExpr.Receiver, ""); err != nil {
		t.Fatal(err)
	}

	if err := testIdentExpr(callExpr.Function, "len"); err != nil {
		t.Fatal(err)
	}
}

func TestParseForDir(t *testing.T) {
	t.Run("regular @for", func(t *testing.T) {
		inp := "@for(i = 0; i < 10; i++){{ i }}@end"

		forDir, err := parseDirective[*ast.ForDir](inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}

		if err := testToken(forDir, token.FOR); err != nil {
			t.Fatal(err)
		}

		if err = testTokPosition(forDir.Pos(), &position.Pos{EndCol: 34}); err != nil {
			t.Fatal(err)
		}

		err = testTokPosition(forDir.Block.Pos(), &position.Pos{StartCol: 24, EndCol: 30})
		if err != nil {
			t.Fatal(err)
		}

		if forDir.Init.String() != `(i = 0)` {
			t.Fatalf("forDir.Init.String() is not '(i = 0)', got %s", forDir.Init)
		}

		if forDir.Cond.String() != `(i < 10)` {
			t.Fatalf("forDir.Cond.String() is not '(i < 10)', got %s", forDir.Cond)
		}

		if forDir.Post.String() != `(i++)` {
			t.Fatalf("forDir.Post.String() is not '(i++)', got %s", forDir.Post)
		}

		actual := strings.TrimSpace(forDir.Block.String())
		if actual != "{{ i }}" {
			t.Fatalf("actual is not '%q', got %q", "{{ i }}", actual)
		}

		if forDir.ElseBlock != nil {
			t.Fatalf("forDir.ElseBlock is not nil, got %T", forDir.ElseBlock)
		}
	})

	t.Run("@for with @else block", func(t *testing.T) {
		inp := `@for(i = 0; i < 0; i++){{ i }}@elseEmpty@end`

		forDir, err := parseDirective[*ast.ForDir](inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}

		if err := testToken(forDir, token.FOR); err != nil {
			t.Fatal(err)
		}

		if forDir.ElseBlock == nil {
			t.Fatalf("forDir.ElseBlock is nil")
		}

		if forDir.ElseBlock.String() != "Empty" {
			t.Fatalf("forDir.ElseBlock.String() is not 'Empty', got %s", forDir.ElseBlock)
		}
	})

	t.Run("infinite @for", func(t *testing.T) {
		inp := `@for(;;)1@end`

		forDir, err := parseDirective[*ast.ForDir](inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}

		if err := testToken(forDir, token.FOR); err != nil {
			t.Fatal(err)
		}

		if _, ok := forDir.Init.(*ast.Empty); !ok {
			t.Fatalf("forDir.Init is not *ast.Empty, got %T", forDir.Init)
		}

		if _, ok := forDir.Cond.(*ast.Empty); !ok {
			t.Fatalf("forDir.Cond is not *ast.Empty, got %T", forDir.Cond)
		}

		if _, ok := forDir.Post.(*ast.Empty); !ok {
			t.Fatalf("forDir.Post is not *ast.Empty, got %T", forDir.Post)
		}

		if forDir.Block.String() != "1" {
			t.Fatalf("forDir.Block.String() is not '1', got %s", forDir.Block)
		}
	})
}

func TestParseEachDir(t *testing.T) {
	t.Run("regular @each", func(t *testing.T) {
		inp := "@each(name in ['anna', 'serhii']){{ name }}@end"

		eachDir, err := parseDirective[*ast.EachDir](inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}

		err = testTokPosition(eachDir.Pos(), &position.Pos{EndCol: 46})
		if err != nil {
			t.Fatal(err)
		}
		if err := testToken(eachDir, token.EACH); err != nil {
			t.Fatal(err)
		}

		err = testTokPosition(eachDir.Block.Pos(), &position.Pos{
			StartCol: 33,
			EndCol:   42,
		})
		if err != nil {
			t.Fatal(err)
		}

		if eachDir.Var.String() != `name` {
			t.Fatalf("eachDir.Var.String() is not 'name', got %s", eachDir.Var)
		}

		if eachDir.Arr.String() != `["anna", "serhii"]` {
			t.Fatalf(`eachDir.Arr.String() is not '["anna", "serhii"]', got %s`, eachDir.Arr)
		}

		actual := eachDir.Block.String()
		if actual != "{{ name }}" {
			t.Fatalf("actual is not %q, got %q", "{{ name }}", actual)
		}

		if eachDir.ElseBlock != nil {
			t.Fatalf("eachDir.ElseBlock is not nil, got %T", eachDir.ElseBlock)
		}
	})

	t.Run("@each with @else", func(t *testing.T) {
		inp := `@each(v in []){{ v }}@elseTest@end`

		eachDir, err := parseDirective[*ast.EachDir](inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}

		if err := testToken(eachDir, token.EACH); err != nil {
			t.Fatal(err)
		}

		if eachDir.ElseBlock.String() != "Test" {
			t.Fatalf("eachDir.ElseBlock.String() is not 'Test', got %s", eachDir.ElseBlock)
		}
	})
}

func TestParseEmptyBlock(t *testing.T) {
	cases := []struct {
		id        uint
		inp       string
		endColPos uint
		tok       token.TokenType
	}{
		{10, "@if(x)@end", 9, token.IF},
		{11, "@if(x)    @end", 13, token.IF},
		{20, "@if(x)1@else@end", 15, token.IF},
		{30, "@if(x)@else@end", 14, token.IF},
		{40, "@if(x)@else1@end", 15, token.IF},
		{41, "@if(x) @else @end", 16, token.IF},
		{50, "@each(x in a)@end", 16, token.EACH},
		{51, "@each(x in a)  @end", 18, token.EACH},
		{60, "@each(x in a)@else@end", 21, token.EACH},
		{70, "@each(x in a)1@else@end", 22, token.EACH},
		{80, "@each(x in a)@else1@end", 22, token.EACH},
		{81, "@each(x in a) @else @end", 23, token.EACH},
		{90, "@for(i = 0; i < x; i++)@end", 26, token.FOR},
		{100, "@for(i = 0; i < x; i++)@else@end", 31, token.FOR},
		{110, "@for(i = 0; i < x; i++)1@else@end", 32, token.FOR},
		{120, "@for(i = 0; i < x; i++)@else1@end", 32, token.FOR},
		{130, "@for(i = 0; i < x; i++) @else @end", 33, token.FOR},
		{140, "@insert('x')@end", 15, token.INSERT},
		{150, "@component('x')@end", 18, token.COMPONENT},
		{160, "@component('x')@pass('x')@end@end", 32, token.COMPONENT},
		{170, "@component('x')@passif(x, 'x')@end@end", 37, token.COMPONENT},
		{
			180,
			"@component('x')@passif(x, 'x')@end@passif(y, 'y')@end@end",
			56,
			token.COMPONENT,
		},
		{190, "@component('x')@pass('x')@end@pass('y')@end@end", 46, token.COMPONENT},
		{
			200,
			"@component('x')@pass('x')@end@pass('y')@end@passif(x, 'name')@end@end",
			68,
			token.COMPONENT,
		},
		{210, "@if(x)1@elseif(y)@end", 20, token.IF},
		{220, "@if(x)@elseif(y)@end", 19, token.IF},
		{230, "@if(x)@elseif(y)1@else@end", 25, token.IF},
		{240, "@if(x)@elseif(y)@elseif(z)@end", 29, token.IF},
		{250, "@if(x)@elseif(y)@elseif(z)1@else@end", 35, token.IF},
		{260, "@if(x) @elseif(y) @elseif(z) @else @end", 38, token.IF},
	}

	for _, tc := range cases {
		chunks, err := parseChunks(tc.inp, defaultParseOpts)
		if err != nil {
			t.Fatalf("Case: %d. %v", tc.id, err)
		}

		node, ok := chunks[0].(ast.NodeWithChunks)
		if !ok {
			t.Fatalf("Case: %d. chunks[0] is not a NodeWithChunks, got %T", tc.id, chunks[0])
		}

		if err := testToken(node, tc.tok); err != nil {
			t.Fatalf("Case: %d. %v", tc.id, err)
		}

		if err = testTokPosition(node.Pos(), &position.Pos{EndCol: tc.endColPos}); err != nil {
			t.Fatalf("Case: %d. %v", tc.id, err)
		}

		if len(node.AllChunks()) != 0 {
			t.Fatalf("len(chunk.AllChunks()) has to be empty, got %d", len(node.AllChunks()))
		}
	}
}

func TestParseObjExpr(t *testing.T) {
	inp := `{{ {"father": {name: "John"},} }}`

	objExpr, err := parseEmbedded[*ast.ObjExpr](inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	err = testTokPosition(objExpr.Pos(), &position.Pos{
		StartCol: 3,
		EndCol:   29,
	})

	if err != nil {
		t.Fatal(err)
	}

	if len(objExpr.Pairs) != 1 {
		t.Fatalf("len(objExpr.Pairs) is not 1, got %d", len(objExpr.Pairs))
	}

	if objExpr.String() != `{"father": {"name": "John"}}` {
		t.Fatalf(`objExpr.String() is not '{"father": {"name": "John"}}', got %s`, objExpr)
	}

	nested, ok := objExpr.Pairs["father"].(*ast.ObjExpr)
	if !ok {
		t.Fatalf("objExpr.Pairs['father'] is not a ObjExpr, got %T", objExpr.Pairs["father"])
	}

	if err := testStrExpr(nested.Pairs["name"], "John"); err != nil {
		t.Fatal(err)
	}

	if err := testToken(objExpr, token.LBRACE); err != nil {
		t.Fatal(err)
	}
}

func TestParseObjWithShorthandKeyNotation(t *testing.T) {
	inp := `{{ { name, age } }}`

	objExpr, err := parseEmbedded[*ast.ObjExpr](inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	if err := testToken(objExpr, token.LBRACE); err != nil {
		t.Fatal(err)
	}

	err = testTokPosition(objExpr.Pos(), &position.Pos{
		StartCol: 3,
		EndCol:   15,
	})

	if err != nil {
		t.Fatal(err)
	}

	if len(objExpr.Pairs) != 2 {
		t.Fatalf("len(objExpr.Pairs) is not 2, got %d", len(objExpr.Pairs))
	}
}

func TestParseText(t *testing.T) {
	inp := "<div><span>Hello</span></div>"

	text, err := parseDirective[*ast.Text](inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	if err := testToken(text, token.TEXT); err != nil {
		t.Fatal(err)
	}

	err = testTokPosition(text.Pos(), &position.Pos{
		StartCol: 0,
		EndCol:   28,
	})

	if err != nil {
		t.Fatal(err)
	}
}

func TestParseDotExp(t *testing.T) {
	inp := "{{ person.father.name }}"

	dotExpr, err := parseEmbedded[*ast.DotExpr](inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	if err := testToken(dotExpr, token.IDENT); err != nil {
		t.Fatal(err)
	}

	// position of the last dot between "father" and "name"
	err = testTokPosition(dotExpr.Pos(), &position.Pos{
		StartCol: 3,
		EndCol:   20,
	})

	if err != nil {
		t.Fatal(err)
	}

	if dotExpr.String() != "((person.father).name)" {
		t.Fatalf("dotExpr.String() is not '((person.father).name)', got %s", dotExpr)
	}

	if err := testIdentExpr(dotExpr.Key, "name"); err != nil {
		t.Fatal(err)
	}

	if dotExpr.Left == nil {
		t.Fatalf("dotExpr.Left is nil")
	}

	leftDotExpr, ok := dotExpr.Left.(*ast.DotExpr)
	if leftDotExpr == nil {
		t.Fatalf("leftDotExpr is nil")
		return
	}

	if !ok {
		t.Fatalf("dotExpr.Left is not a DotExpr, got %T", dotExpr.Left)
	}

	if err := testIdentExpr(leftDotExpr.Key, "father"); err != nil {
		t.Fatal(err)
	}

	if err := testIdentExpr(leftDotExpr.Left, "person"); err != nil {
		t.Fatal(err)
	}
}

func TestParseBreakDir(t *testing.T) {
	inp := `@break`

	breakDir, err := parseDirective[*ast.BreakDir](inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	if err := testToken(breakDir, token.BREAK); err != nil {
		t.Fatal(err)
	}
}

func TestParseContinueDir(t *testing.T) {
	inp := `@continue`

	continueDir, err := parseDirective[*ast.ContinueDir](inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	if err := testToken(continueDir, token.CONTINUE); err != nil {
		t.Fatal(err)
	}
}

func TestParseBreakifDir(t *testing.T) {
	inp := `@breakif(true)`

	breakifDir, err := parseDirective[*ast.BreakifDir](inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	err = testTokPosition(breakifDir.Pos(), &position.Pos{
		StartCol: 0,
		EndCol:   13,
	})

	if err != nil {
		t.Fatal(err)
	}

	if err := testToken(breakifDir, token.BREAKIF); err != nil {
		t.Fatal(err)
	}

	if err := testBoolExpr(breakifDir.Cond, true); err != nil {
		t.Fatal(err)
	}

	expect := "@breakif(true)"
	if breakifDir.String() != expect {
		t.Fatalf("breakifDir.String() is not '%s', got %s", expect, breakifDir)
	}
}

func TestParseContinueifDir(t *testing.T) {
	inp := "@continueif(false)"

	continueifDir, err := parseDirective[*ast.ContinueifDir](inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	err = testTokPosition(continueifDir.Pos(), &position.Pos{
		StartCol: 0,
		EndCol:   17,
	})

	if err != nil {
		t.Fatal(err)
	}

	if err := testToken(continueifDir, token.CONTINUEIF); err != nil {
		t.Fatal(err)
	}

	if err := testBoolExpr(continueifDir.Cond, false); err != nil {
		t.Fatal(err)
	}

	expect := "@continueif(false)"
	if continueifDir.String() != expect {
		t.Fatalf("continueifDir.String() is not '%s', got %s", expect, continueifDir)
	}
}

func TestParseCompDir(t *testing.T) {
	t.Run("@component without @pass", func(t *testing.T) {
		inp := "<ul>@component('components/book-card', { c: card })@end</ul>"
		chunks, err := parseChunks(inp, parseOpts{chunksCount: 3, checkErrors: true})
		if err != nil {
			t.Fatal(err)
		}

		compDir, ok := chunks[1].(*ast.CompDir)
		if !ok {
			t.Fatalf("chunks[1] is not a ComponentDir, got %T", chunks[1])
		}

		err = testTokPosition(compDir.Pos(), &position.Pos{
			StartCol: 4,
			EndCol:   54,
		})

		if err != nil {
			t.Fatal(err)
		}

		if err := testToken(compDir, token.COMPONENT); err != nil {
			t.Fatal(err)
		}

		if err := testStrExpr(compDir.Name, "components/book-card"); err != nil {
			t.Fatal(err)
		}

		if len(compDir.Argument.Pairs) != 1 {
			t.Fatalf("len(compDir.Argument.Pairs) is not 1, got %d", len(compDir.Argument.Pairs))
		}

		if err := testIdentExpr(compDir.Argument.Pairs["c"], "card"); err != nil {
			t.Fatal(err)
		}

		if len(compDir.Passes) != 0 {
			t.Fatalf("len(compDir.Passes) is not 0, got %d", len(compDir.Passes))
		}

		if compDir.DefaultPass != nil {
			t.Fatalf("compDir.DefaultPass is not nil, got %v", compDir.DefaultPass)
		}

		expect := `@component("components/book-card", {"c": card})@end`
		if compDir.String() != expect {
			t.Fatalf(`compDir.String() is not '%s', got %s`, expect, compDir)
		}
	})

	t.Run("@component with default @pass", func(t *testing.T) {
		inp := `@component("components/book-card")<h1>Header</h1>@end`

		compDir, err := parseDirective[*ast.CompDir](inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}

		if err := testToken(compDir, token.COMPONENT); err != nil {
			t.Fatal(err)
		}

		if len(compDir.Passes) != 0 {
			t.Fatalf("len(compDir.Passes) is not 0, got %d", len(compDir.Passes))
		}

		passDir := compDir.DefaultPass
		name := passDir.Name.Val
		if name != "" {
			t.Fatalf("passDir.Name.Val must be empty string, got: %s", name)
		}

		expect := `<h1>Header</h1>`
		if passDir.Block.String() != expect {
			t.Fatalf(
				"passDir.Block.String() is not '%q', got '%q'",
				expect,
				passDir.Block,
			)
		}
	})

	t.Run("@component with 2 @pass directives", func(t *testing.T) {
		inp := `<ul>
			@component("components/book-card", { c: card })
				@pass("header")<h1>Header</h1>@end
				@pass("footer")<footer>Footer</footer>@end
			@end
		</ul>`

		chunks, err := parseChunks(inp, parseOpts{chunksCount: 3, checkErrors: true})
		if err != nil {
			t.Fatal(err)
		}

		compDir, ok := chunks[1].(*ast.CompDir)
		if !ok {
			t.Fatalf("chunks[1] is not a ComponentDir, got %T", chunks[1])
		}

		err = testTokPosition(compDir.Pos(), &position.Pos{
			StartLine: 1,
			EndLine:   4,
			StartCol:  3, // because tabs before @component
			EndCol:    6, // because tabs before @end
		})

		if err != nil {
			t.Fatal(err)
		}

		if err := testToken(compDir, token.COMPONENT); err != nil {
			t.Fatal(err)
		}

		if len(compDir.Passes) != 2 {
			t.Fatalf("len(compDir.Passes) is not 2, got %d", len(compDir.Passes))
		}

		if compDir.DefaultPass != nil {
			t.Fatalf("compDir.DefaultPass is not nil, got %s", compDir.DefaultPass)
		}

		passHeader := compDir.Passes[0]
		passFooter := compDir.Passes[1]

		if err := testStrExpr(passHeader.Name, "header"); err != nil {
			t.Fatal(err)
		}

		expect := `<h1>Header</h1>`
		if passHeader.Block.String() != expect {
			t.Fatalf(
				"passHeader.Block.String() is not '%q', got %q",
				expect,
				passHeader.Block,
			)
		}

		if err := testStrExpr(passFooter.Name, "footer"); err != nil {
			t.Fatal(err)
		}

		expect = `<footer>Footer</footer>`
		if passFooter.Block.String() != expect {
			t.Fatalf(
				"passFooter.Block.String() is not '%q', got %q",
				expect,
				passFooter.Block,
			)
		}
	})

	t.Run("@component with 1 @passif", func(t *testing.T) {
		inp := `@component('user')@passif(false, 'name')Test2@end@end`
		compDir, err := parseDirective[*ast.CompDir](inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}

		err = testTokPosition(compDir.Pos(), &position.Pos{EndCol: 52})
		if err != nil {
			t.Fatal(err)
		}

		if err := testToken(compDir, token.COMPONENT); err != nil {
			t.Fatal(err)
		}

		passDir := compDir.Passes[0]

		if err := testBoolExpr(passDir.Cond, false); err != nil {
			t.Fatal(err)
		}

		if passDir.Name.Val != "name" {
			t.Fatalf("passDir.Name.Val is not 'name', got %s", passDir.Name)
		}

		block := passDir.Block.String()
		if block != "Test2" {
			t.Fatalf("passDir.Block.String() is not 'Test2', got %s", block)
		}

		expect := `@passif(false, "name")Test2@end`
		if passDir.String() != expect {
			t.Fatalf("passDir.String() is not '%s', got %s", expect, passDir)
		}
	})
}

func TestParseSlotDir(t *testing.T) {
	t.Run("named slot", func(t *testing.T) {
		inp := "<h2>@slot('header')</h2>"
		chunks, err := parseChunks(inp, parseOpts{chunksCount: 3, checkErrors: true})
		if err != nil {
			t.Fatal(err)
		}

		slotDir, ok := chunks[1].(*ast.SlotDir)
		if !ok {
			t.Fatalf("chunks[1] is not a SlotDir, got %T", chunks[1])
		}

		err = testTokPosition(slotDir.Pos(), &position.Pos{
			StartCol: 4,
			EndCol:   18,
		})

		if err != nil {
			t.Fatal(err)
		}

		if err := testToken(slotDir, token.SLOT); err != nil {
			t.Fatal(err)
		}

		if err := testStrExpr(slotDir.Name, "header"); err != nil {
			t.Fatal(err)
		}

		expect := `@slot("header")`
		if slotDir.String() != expect {
			t.Fatalf("slotDir.String() is not `%s`, got `%s`", expect, slotDir)
		}
	})

	t.Run("default slot without end", func(t *testing.T) {
		inp := `<header>@slot</header>`
		chunks, err := parseChunks(inp, parseOpts{chunksCount: 3, checkErrors: true})
		if err != nil {
			t.Fatal(err)
		}

		slotDir, ok := chunks[1].(*ast.SlotDir)
		if !ok {
			t.Fatalf("chunks[1] is not a SlotDir, got %T", chunks[1])
		}

		if err := testToken(slotDir, token.SLOT); err != nil {
			t.Fatal(err)
		}

		if err := testStrExpr(slotDir.Name, ""); err != nil {
			t.Fatal(err)
		}

		err = testTokPosition(slotDir.Pos(), &position.Pos{
			StartCol: 8,
			EndCol:   12,
		})

		if err != nil {
			t.Fatal(err)
		}

		if slotDir.String() != "@slot" {
			t.Fatalf("slotDir.String() is not @slot, got `%s`", slotDir)
		}
	})
}

func TestParseDumpDir(t *testing.T) {
	inp := `@dump("test", 1 + 2, false)`

	dumpDir, err := parseDirective[*ast.DumpDir](inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	if len(dumpDir.Args) != 3 {
		t.Fatalf("len(dumpDir.Args) is not 3, got %d", len(dumpDir.Args))
	}

	err = testTokPosition(dumpDir.Pos(), &position.Pos{
		StartCol: 0,
		EndCol:   26,
	})

	if err != nil {
		t.Fatal(err)
	}

	if err := testToken(dumpDir, token.DUMP); err != nil {
		t.Fatal(err)
	}

	if err := testStrExpr(dumpDir.Args[0], "test"); err != nil {
		t.Fatal(err)
	}

	if err := testInfixExpr(dumpDir.Args[1], 1, "+", 2); err != nil {
		t.Fatal(err)
	}

	if err := testBoolExpr(dumpDir.Args[2], false); err != nil {
		t.Fatal(err)
	}
}
