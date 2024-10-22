package evaluator

import "testing"

func TestEvalIntFunctions(t *testing.T) {
	tests := []struct {
		inp      string
		expected string
	}{
		// float func
		{`{{ 1.float() }}`, "1.0"},
		{`{{ 321.float().int() }}`, "321"},
		// abs func
		{`{{ 1.abs() }}`, "1"},
		{`{{ -1.abs() }}`, "1"},
		{`{{ 0.abs() }}`, "0"},
		{`{{ -999999.abs() }}`, "999999"},
		{`{{ 999999.abs() }}`, "999999"},
		// str func
		{`{{ 1.str() }}`, "1"},
		{`{{ 321.str() }}`, "321"},
		{`{{ -1.str() }}`, "-1"},
		{`{{ -9999999999.str() }}`, "-9999999999"},
		{`{{ 9999999999.str() }}`, "9999999999"},
	}

	for _, tc := range tests {
		evaluationExpected(t, tc.inp, tc.expected)
	}
}
