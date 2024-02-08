package evaluator

import (
	"testing"

	"github.com/textwire/textwire/object"
)

func TestIsTruthy(t *testing.T) {
	tests := []struct {
		inp      object.Object
		expected bool
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

	for _, tt := range tests {
		result := isTruthy(tt.inp)

		if result != tt.expected {
			t.Errorf("isTruthy(%v) returned %t, expected %t", tt.inp, result, tt.expected)
		}
	}
}

func TestNativeToBooleanObject(t *testing.T) {
	tests := []struct {
		inp      bool
		expected object.Object
	}{
		{true, TRUE},
		{false, FALSE},
	}

	for _, tt := range tests {
		result := nativeBoolToBooleanObject(tt.inp)

		if result != tt.expected {
			t.Errorf("nativeBoolToBooleanObject(%t) returned %s, expected %s", tt.inp, result, tt.expected)
		}
	}
}
