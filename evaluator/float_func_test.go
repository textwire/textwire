package evaluator

import "testing"

func TestEvalFloatFunctions(t *testing.T) {
	tests := []struct {
		inp      string
		expected string
	}{
		// int
		{`{{ 13.999.int() }}`, "13"},
		{`{{ 133.999124.int() }}`, "133"},
		{`{{ 1234567.999124.int() }}`, "1234567"},
		{`{{ -133.999124.int() }}`, "-133"},
		// str
		{`{{ 13.999.str() }}`, "13.999"},
		{`{{ 133.999124.str() }}`, "133.999124"},
		{`{{ 1234567890.1234567.str() }}`, "1234567890.1234567"},
		// abs
		{`{{ 1.0.abs() }}`, "1.0"},
		{`{{ -1.0.abs() }}`, "1.0"},
		{`{{ 0.0.abs() }}`, "0.0"},
		{`{{ -999999.55.abs() }}`, "999999.55"},
		{`{{ 999999.55.abs() }}`, "999999.55"},
		// ceil
		{`{{ 1.0.ceil() }}`, "1"},
		{`{{ 1.1.ceil() }}`, "2"},
		{`{{ 1.1.ceil() }}`, "2"},
		{`{{ 5.125.ceil() }}`, "6"},
		{`{{ 1.9.ceil() }}`, "2"},
		{`{{ 0.0.ceil() }}`, "0"},
		{`{{ -1.0.ceil() }}`, "-1"},
		{`{{ -1.1.ceil() }}`, "-1"},
		{`{{ -1.543.ceil() }}`, "-1"},
		{`{{ -1.9.ceil() }}`, "-1"},
		// floor
		{`{{ 1.0.floor() }}`, "1"},
		{`{{ 1.1.floor() }}`, "1"},
		{`{{ 1.9.floor() }}`, "1"},
		{`{{ 5.125.floor() }}`, "5"},
		{`{{ 0.0.floor() }}`, "0"},
		{`{{ -1.0.floor() }}`, "-1"},
		{`{{ -1.1.floor() }}`, "-2"},
		{`{{ -1.9.floor() }}`, "-2"},
		{`{{ -5.125.floor() }}`, "-6"},
		// round
		{`{{ 1.0.round() }}`, "1"},
		{`{{ 1.1.round() }}`, "1"},
		{`{{ 1.4.round() }}`, "1"},
		{`{{ 1.4999.round() }}`, "1"},
		{`{{ 1.5.round() }}`, "2"},
		{`{{ 1.6.round() }}`, "2"},
		{`{{ 1.9.round() }}`, "2"},
		{`{{ 5.125.round() }}`, "5"},
		{`{{ 0.0.round() }}`, "0"},
		{`{{ -1.0.round() }}`, "-1"},
		{`{{ -1.1.round() }}`, "-1"},
		{`{{ -1.9.round() }}`, "-2"},
		{`{{ -5.125.round() }}`, "-5"},
	}

	for _, tc := range tests {
		evaluationExpected(t, tc.inp, tc.expected)
	}
}
