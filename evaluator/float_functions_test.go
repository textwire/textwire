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
		// abs func
		{`{{ 1.0.abs() }}`, "1.0"},
		{`{{ -1.0.abs() }}`, "1.0"},
		{`{{ 0.0.abs() }}`, "0.0"},
		{`{{ -999999.55.abs() }}`, "999999.55"},
		{`{{ 999999.55.abs() }}`, "999999.55"},
	}

	for _, tc := range tests {
		evaluationExpected(t, tc.inp, tc.expected)
	}
}
