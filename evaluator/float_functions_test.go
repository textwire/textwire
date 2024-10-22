package evaluator

import "testing"

func TestEvalFloatFunctions(t *testing.T) {
	tests := []struct {
		inp      string
		expected string
	}{
		// int func
		{`{{ 13.999.int() }}`, "13"},
		{`{{ 133.999124.int() }}`, "133"},
		{`{{ 1234567.999124.int() }}`, "1234567"},
		// str func
		{`{{ 13.999.str() }}`, "13.999"},
		{`{{ 133.999124.str() }}`, "133.999124"},
		{`{{ 1234567890.1234567.str() }}`, "1234567890.1234567"},
	}

	for _, tc := range tests {
		evaluationExpected(t, tc.inp, tc.expected)
	}
}
