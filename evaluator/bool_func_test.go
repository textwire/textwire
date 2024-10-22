package evaluator

import "testing"

func TestEvalBinaryFunctions(t *testing.T) {
	tests := []struct {
		inp      string
		expected string
	}{
		// binary
		{`{{ true.binary() }}`, "1"},
		{`{{ false.binary() }}`, "0"},
		{`{{ !true.binary() }}`, "0"},
		{`{{ !false.binary() }}`, "1"},
		{`{{ true.binary().float() }}`, "1.0"},
		{`{{ false.binary().float() }}`, "0.0"},
	}

	for _, tc := range tests {
		evaluationExpected(t, tc.inp, tc.expected)
	}
}
