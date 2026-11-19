package evaluator

import "testing"

func TestEvalIntFunctions(t *testing.T) {
	cases := []struct {
		name   string
		inp    string
		expect string
	}{
		// float
		{"int to float", `{{ 1.float() }}`, "1.0"},
		{"int to float and back to int", `{{ 321.float().int() }}`, "321"},
		// abs
		{"abs of positive int", `{{ 1.abs() }}`, "1"},
		{"abs on literal before negation", `{{ -1.abs() }}`, "-1"},
		{"abs of parenthesized negative int", `{{ (-1).abs() }}`, "1"},
		{"abs of zero", `{{ 0.abs() }}`, "0"},
		{"abs of negative literal negated after abs", `{{ -999999.abs() }}`, "-999999"},
		{"abs of large positive int", `{{ 999999.abs() }}`, "999999"},
		// str
		{"int to string", `{{ 1.str() }}`, "1"},
		{"large int to string", `{{ 321.str() }}`, "321"},
		{"negative parenthesized int to string", `{{ (-1).str() }}`, "-1"},
		{"large negative int to string", `{{ (-9999999999).str() }}`, "-9999999999"},
		{"large positive int to string", `{{ 9999999999.str() }}`, "9999999999"},
		// len
		{"len of single digit", `{{ 1.len() }}`, "1"},
		{"len of three digits", `{{ 321.len() }}`, "3"},
		{"len of negative literal applies negation after len", `{{ -1.len() }}`, "-1"},
		{"len of zero", `{{ 0.len() }}`, "1"},
		{"len of negative literal with trailing negation", `{{ -9999999999.len() }}`, "-10"},
		{"len of parenthesized negative int", `{{ (-9999999999).len() }}`, "10"},
		{"len of ten digit int", `{{ 9999999999.len() }}`, "10"},
		// decimal
		{"decimal of zero", `{{ 0.decimal() }}`, "0.00"},
		{"decimal of hundred", `{{ 100.decimal() }}`, "100.00"},
		{"decimal of 2352", `{{ 2352.decimal() }}`, "2352.00"},
		{"decimal with custom separator", `{{ 1000.decimal('_') }}`, "1000_00"},
		{"decimal with separator and precision", `{{ 9000.decimal('_', 10) }}`, "9000_0000000000"},
		{"decimal with zero precision", `{{ 100.decimal('|', 0) }}`, "100"},
		{"decimal with one precision", `{{ 100.decimal('|', 1) }}`, "100|0"},
		{
			"decimal of parenthesized negative with zero precision",
			`{{ (-100).decimal('|', 0) }}`,
			"-100",
		},
		{"decimal with comma separator", `{{ 1.decimal(',') }}`, "1,00"},
		{"decimal with negative precision", `{{ 1.decimal('.', -1) }}`, "1"},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, nil, tc.name)
	}
}
