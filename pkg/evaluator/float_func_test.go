package evaluator

import "testing"

func TestEvalFloatFunctions(t *testing.T) {
	cases := []struct {
		name   string
		inp    string
		expect string
	}{
		// int
		{"float to int", `{{ 13.999.int() }}`, "13"},
		{"float to int truncates fraction", `{{ 133.999124.int() }}`, "133"},
		{"large float to int truncates fraction", `{{ 1234567.999124.int() }}`, "1234567"},
		{"negative float to int truncates fraction", `{{ -133.999124.int() }}`, "-133"},
		{
			"float at int64 boundary to int",
			`{{ (-9223372036854775808.0).int() }}`,
			"-9223372036854775808",
		},
		// str
		{"float to string", `{{ 13.999.str() }}`, "13.999"},
		{"float to string keeps fraction", `{{ 133.999124.str() }}`, "133.999124"},
		{"large float to string", `{{ 1234567890.1234567.str() }}`, "1234567890.1234567"},
		// abs
		{"abs of positive float", `{{ 1.0.abs() }}`, "1.0"},
		{"abs on literal before negation", `{{ -1.0.abs() }}`, "-1.0"},
		{"abs of zero", `{{ 0.0.abs() }}`, "0.0"},
		{"abs of negative literal negated after abs", `{{ -999999.55.abs() }}`, "-999999.55"},
		{"abs of positive large float", `{{ 999999.55.abs() }}`, "999999.55"},
		// ceil
		{"ceil of whole number", `{{ 1.0.ceil() }}`, "1"},
		{"ceil of small fraction", `{{ 1.1.ceil() }}`, "2"},
		{"ceil of 1.1", `{{ 1.1.ceil() }}`, "2"},
		{"ceil of 5.125", `{{ 5.125.ceil() }}`, "6"},
		{"ceil of 1.9", `{{ 1.9.ceil() }}`, "2"},
		{"ceil of zero", `{{ 0.0.ceil() }}`, "0"},
		{"ceil of negative whole number", `{{ -1.0.ceil() }}`, "-1"},
		{"ceil of negative fraction", `{{ -1.1.ceil() }}`, "-2"},
		{"ceil of negative 1.543", `{{ -1.543.ceil() }}`, "-2"},
		{"ceil of negative 1.9", `{{ -1.9.ceil() }}`, "-2"},
		// floor
		{"floor of whole number", `{{ 1.0.floor() }}`, "1"},
		{"floor of 1.1", `{{ 1.1.floor() }}`, "1"},
		{"floor of 1.9", `{{ 1.9.floor() }}`, "1"},
		{"floor of 5.125", `{{ 5.125.floor() }}`, "5"},
		{"floor of zero", `{{ 0.0.floor() }}`, "0"},
		{"floor of negative whole number", `{{ -1.0.floor() }}`, "-1"},
		{"floor of negative 1.1", `{{ -1.1.floor() }}`, "-1"},
		{"floor of negative 1.9", `{{ -1.9.floor() }}`, "-1"},
		{"floor of negative 5.125", `{{ -5.125.floor() }}`, "-5"},
		// round
		{"round of whole number", `{{ 1.0.round() }}`, "1"},
		{"round of 1.1", `{{ 1.1.round() }}`, "1"},
		{"round 1.4 down", `{{ 1.4.round() }}`, "1"},
		{"round 1.4999 down", `{{ 1.4999.round() }}`, "1"},
		{"round 1.5 up", `{{ 1.5.round() }}`, "2"},
		{"round 1.6 up", `{{ 1.6.round() }}`, "2"},
		{"round 1.9 up", `{{ 1.9.round() }}`, "2"},
		{"round 5.125 down", `{{ 5.125.round() }}`, "5"},
		{"round zero", `{{ 0.0.round() }}`, "0"},
		{"round negative whole number", `{{ -1.0.round() }}`, "-1"},
		{"round negative 1.1 toward zero", `{{ -1.1.round() }}`, "-1"},
		{"round negative 1.9 away from zero", `{{ -1.9.round() }}`, "-2"},
		{"round negative 5.125 toward zero", `{{ -5.125.round() }}`, "-5"},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, nil, tc.name)
	}
}
