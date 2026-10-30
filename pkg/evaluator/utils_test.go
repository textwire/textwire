package evaluator

import (
	"testing"

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
