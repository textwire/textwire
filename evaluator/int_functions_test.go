package evaluator

import "testing"

func TestEvalIntFunctions(t *testing.T) {
	tests := []struct {
		inp      string
		expected string
	}{
		{`{{ 1.float() }}`, "1.0"},
		{`{{ 321.float().int() }}`, "321"},
	}

	for _, tt := range tests {
		evaluationExpected(t, tt.inp, tt.expected)
	}
}
