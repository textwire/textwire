package evaluator

import (
	"reflect"
	"strings"
	"testing"

	"codeberg.org/textwire/textwire/v5/pkg/ast"
	"codeberg.org/textwire/textwire/v5/pkg/token"
	"codeberg.org/textwire/textwire/v5/pkg/value"
)

func TestIsTruthy(t *testing.T) {
	cases := []struct {
		inp    value.Literal
		expect bool
	}{
		{nil, false},
		{NIL, false},
		{TRUE, true},
		{FALSE, false},
		{&value.Int{Val: 0}, false},
		{&value.Int{Val: 1}, true},
		{&value.Int{Val: -1}, true},
		{&value.Float{Val: 0.0}, false},
		{&value.Float{Val: 1.0}, true},
		{&value.Float{Val: -1.0}, true},
		{&value.Str{Val: ""}, false},
		{&value.Str{Val: "x"}, true},
		{&value.Str{Val: "anna"}, true},
		{&value.Arr{Elements: nil}, false},
		{&value.Obj{Pairs: nil}, false},
	}

	for _, tc := range cases {
		result := isTruthy(tc.inp)

		if result != tc.expect {
			t.Errorf("isTruthy(%v) returned %t, expect %t", tc.inp, result, tc.expect)
		}
	}
}

func TestNativeBoolToBoolObj(t *testing.T) {
	cases := []struct {
		inp    bool
		expect value.Value
	}{
		{true, TRUE},
		{false, FALSE},
	}

	for _, tc := range cases {
		result := nativeBoolToBoolObj(tc.inp)

		if result != tc.expect {
			t.Errorf(
				"nativeBoolToBoolObj(%t) returned %s, expect %s",
				tc.inp,
				result,
				tc.expect,
			)
		}
	}
}

func TestStrIsInt(t *testing.T) {
	tc := []struct {
		name   string
		inp    string
		expect bool
	}{
		{"Non-integer string", "anna", false},
		{"Positive integer", "123", true},
		{"Negative integer", "-123", true},
		{"Zero as integer", "0", true},
		{"Negative one", "-1", true},
		{"Decimal number with fraction", "123.23", false},
		{"Decimal number ending with zero", "123.0", false},
		{"Very large integer", "12345678901234567890", true},
		{"Positive sign", "+123", true},
		{"Empty string", "", false},
		{"Only sign", "-", false},
		{"Leading space", " 1", false},
		{"Exponent notation", "1e5", false},
		{"Special float word", "nan", false},
		{"Two signs", "+-33", false},
		{"Tree signs", "+-+21", false},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			got := strIsInt(tt.inp)
			if got != tt.expect {
				t.Errorf("expect %v, got %v", tt.expect, got)
			}
		})
	}
}

func TestIsValidFloat(t *testing.T) {
	tc := []struct {
		name   string
		inp    string
		expect bool
	}{
		{"Simple decimal", "1.5", true},
		{"Negative decimal", "-0.5", true},
		{"Positive sign decimal", "+1.5", true},
		{"Trailing dot", "1.", true},
		{"Large decimal", "123456789012345678901.5", true},
		{"Integer without dot", "123", false},
		{"Empty string", "", false},
		{"Only sign", "-", false},
		{"Only dot", ".", false},
		{"Two dots", "1.2.3", false},
		{"Letters", "1.a", false},
		{"Nan", "nan", false},
		{"Inf", "inf", false},
		{"Infinity", "infinity", false},
		{"Exponent notation", "1.5e3", false},
		{"Leading space", " 1.5", false},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidFloat(tt.inp)

			if got != tt.expect {
				t.Errorf("expect %v, got %v", tt.expect, got)
			}
		})
	}
}

func TestExprToText(t *testing.T) {
	t.Run("nil expression", func(t *testing.T) {
		if got := exprToText(nil); got != "" {
			t.Errorf("expect %q, got %q", "", got)
		}
	})

	t.Run("identifier", func(t *testing.T) {
		expr := ast.NewIdentExpr(token.Token{Lit: "user"}, "user")
		if got := exprToText(expr); got != "user" {
			t.Errorf("expect %q, got %q", "user", got)
		}
	})

	t.Run("integer", func(t *testing.T) {
		expr := ast.NewIntExpr(token.Token{Lit: "1"}, 1)
		if got := exprToText(expr); got != "1" {
			t.Errorf("expect %q, got %q", "1", got)
		}
	})

	t.Run("infix", func(t *testing.T) {
		expr := ast.NewInfixExpr(
			token.Token{Lit: "+"},
			ast.NewIdentExpr(token.Token{Lit: "a"}, "a"),
			"+",
		)

		expr.Right = ast.NewIdentExpr(token.Token{Lit: "b"}, "b")
		if got := exprToText(expr); got != "a + b" {
			t.Errorf("expect %q, got %q", "a + b", got)
		}
	})

	t.Run("index", func(t *testing.T) {
		expr := ast.NewIndexExpr(
			token.Token{Lit: "["},
			ast.NewIdentExpr(token.Token{Lit: "arr"}, "arr"),
		)

		expr.Index = ast.NewIntExpr(token.Token{Lit: "0"}, 0)
		if got := exprToText(expr); got != "arr[0]" {
			t.Errorf("expect %q, got %q", "arr[0]", got)
		}
	})

	t.Run("dot", func(t *testing.T) {
		expr := ast.NewDotExpr(
			token.Token{Lit: "."},
			ast.NewIdentExpr(token.Token{Lit: "user"}, "user"),
		)

		expr.Key = ast.NewIdentExpr(token.Token{Lit: "name"}, "name")
		if got := exprToText(expr); got != "user.name" {
			t.Errorf("expect %q, got %q", "user.name", got)
		}
	})
}

func TestEnsureMinFloatDecimals(t *testing.T) {
	tc := []struct {
		name      string
		inp       string
		separator string
		decimals  int
		expect    string
	}{
		{"Pads missing decimals", "1.5", ".", 2, "1.50"},
		{"Pads with custom separator", "1.5", ",", 2, "1,50"},
		{"Pads with 3 decimals", "1.5", ".", 3, "1.500"},
		{"Enough decimals keeps value", "1.555", ".", 2, "1.555"},
		{"Enough decimals swaps separator", "1.555", ",", 2, "1,555"},
		{"Exact decimals swaps separator", "1.55", ",", 2, "1,55"},
		{"Zero decimals keeps value", "1.5", ",", 0, "1,5"},
		{"Negative value", "-0.5", ",", 2, "-0,50"},
		{"Positive sign", "+1.5", ".", 2, "+1.50"},
		{"Trailing dot gets padded", "1.", ".", 2, "1.00"},
		{"Large value stays exact", "123456789012345678901.5", ".", 2, "123456789012345678901.50"},
		{"No dot returns value as-is", "5", ".", 2, "5"},
		{"Zero fraction padded", "0.0", ".", 3, "0.000"},
		{"Pads to larger decimals", "1.5", ".", 5, "1.50000"},
		{"Multi-character separator", "1.5", "||", 2, "1||50"},
		{"Empty separator", "1.5", "", 2, "150"},
		{"Plus sign kept", "+1.5", ",", 2, "+1,50"},
		{"Trailing dot padded with separator", "1.", ",", 2, "1,00"},
		{"Trailing dot with zero decimals", "1.", ".", 0, "1."},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			got := ensureMinFloatDecimals(tt.inp, tt.separator, tt.decimals)

			if got != tt.expect {
				t.Errorf("expect %q, got %q", tt.expect, got)
			}
		})
	}
}

func TestSplitClusters(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect []string
	}{
		{
			name:   "empty string",
			input:  "",
			expect: nil,
		},
		{
			name:   "single ASCII char",
			input:  "a",
			expect: []string{"a"},
		},
		{
			name:   "plain ASCII",
			input:  "abc",
			expect: []string{"a", "b", "c"},
		},
		{
			name:   "multibyte runes without marks",
			input:  "日本語",
			expect: []string{"日", "本", "語"},
		},
		{
			name:   "precomposed character is a single rune",
			input:  "\u00e9", // é
			expect: []string{"\u00e9"},
		},
		{
			name:   "base plus combining acute (Mn)",
			input:  "e\u0301",
			expect: []string{"e\u0301"},
		},
		{
			name:   "mark in the middle attaches to the preceding char",
			input:  "ae\u0301b",
			expect: []string{"a", "e\u0301", "b"},
		},
		{
			name:   "multiple stacked marks on one base",
			input:  "a\u0301\u0302\u0303",
			expect: []string{"a\u0301\u0302\u0303"},
		},
		{
			name:   "enclosing mark (Me)",
			input:  "a\u20dd", // combining enclosing circle
			expect: []string{"a\u20dd"},
		},
		{
			name:   "spacing combining mark (Mc)",
			input:  "\u0915\u093e", // Devanagari ka + aa vowel sign
			expect: []string{"\u0915\u093e"},
		},
		{
			name:   "mixed Mn, Me, Mc in sequence",
			input:  "a\u0301\u20dd\u093e",
			expect: []string{"a\u0301\u20dd\u093e"},
		},
		{
			name:   "several marked characters in a row",
			input:  "e\u0301a\u0300o\u0308",
			expect: []string{"e\u0301", "a\u0300", "o\u0308"},
		},
		{
			name:   "variation selector (Mn) attaches to its base",
			input:  "\u2764\ufe0f", // red heart with VS16
			expect: []string{"\u2764\ufe0f"},
		},
		{
			name:   "leading combining mark stands alone",
			input:  "\u0301",
			expect: []string{"\u0301"},
		},
		{
			name:   "leading mark followed by base char",
			input:  "\u0301a",
			expect: []string{"\u0301", "a"},
		},
		{
			name:   "multiple leading marks form one cluster",
			input:  "\u0301\u0302a",
			expect: []string{"\u0301\u0302", "a"},
		},
		{
			name:   "space and punctuation are not marks",
			input:  "a b.",
			expect: []string{"a", " ", "b", "."},
		},
		{
			name:   "invalid UTF-8 byte is its own cluster and preserved",
			input:  "a\xffb",
			expect: []string{"a", "\xff", "b"},
		},
		{
			name:   "mark after invalid byte attaches to it",
			input:  "\xff\u0301",
			expect: []string{"\xff\u0301"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := splitClusters(tt.input)

			if !reflect.DeepEqual(got, tt.expect) {
				t.Errorf("splitClusters(%q)\n got: %q\nwant: %q", tt.input, got, tt.expect)
			}

			// Invariant: clusters must reassemble to exactly the input.
			if joined := strings.Join(got, ""); joined != tt.input {
				t.Errorf("clusters do not reassemble: joined %q, input %q", joined, tt.input)
			}
		})
	}
}
