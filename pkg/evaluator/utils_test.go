package evaluator

import (
	"testing"

	"github.com/textwire/textwire/v4/pkg/value"
)

func TestIsTruthy(t *testing.T) {
	cases := []struct {
		inp    value.Value
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
