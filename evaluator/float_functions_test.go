package evaluator

import "testing"

func TestEvalFloatFunctions(t *testing.T) {
	tests := []struct {
		inp      string
		expected string
	}{
		{`{{ 13.999.int() }}`, "13"},
	}

	for _, tt := range tests {
		evaluationExpected(t, tt.inp, tt.expected)
	}
}
