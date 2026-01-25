package evaluator

import (
	"testing"

	"github.com/textwire/textwire/v2/object"
)

func TestIsTruthy(t *testing.T) {
	cases := []struct {
		inp    object.Object
		expect bool
	}{
		{nil, false},
		{NIL, false},
		{TRUE, true},
		{FALSE, false},
		{&object.Int{Value: 0}, false},
		{&object.Int{Value: 1}, true},
		{&object.Int{Value: -1}, true},
		{&object.Float{Value: 0.0}, false},
		{&object.Float{Value: 1.0}, true},
		{&object.Float{Value: -1.0}, true},
		{&object.Str{Value: ""}, false},
		{&object.Str{Value: "x"}, true},
		{&object.Str{Value: "anna"}, true},
		{&object.Array{Elements: nil}, true},
	}

	for _, tc := range cases {
		result := isTruthy(tc.inp)

		if result != tc.expect {
			t.Errorf("isTruthy(%v) returned %t, expect %t", tc.inp, result, tc.expect)
		}
	}
}

func TestNativeToBooleanObject(t *testing.T) {
	cases := []struct {
		inp    bool
		expect object.Object
	}{
		{true, TRUE},
		{false, FALSE},
	}

	for _, tc := range cases {
		result := nativeBoolToBooleanObject(tc.inp)

		if result != tc.expect {
			t.Errorf("nativeBoolToBooleanObject(%t) returned %s, expect %s", tc.inp, result, tc.expect)
		}
	}
}
