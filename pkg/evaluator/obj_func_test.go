package evaluator

import (
	"testing"
)

func TestEvalObjectFunctions(t *testing.T) {
	cases := []struct {
		inp    string
		expect string
	}{
		// json
		{`{{ {}.json() }}`, "{}"},
		{`{{ {one: {two: {}}}.json() }}`, `{"one":{"two":{}}}`},
		{
			inp:    `{{ {name: "Chiori", game: "Genshin Impact"}.json() }}`,
			expect: `{"game":"Genshin Impact","name":"Chiori"}`,
		},
		{
			inp:    `{{ user = {address: {street: "Via Emilio Morosini", city: "Rome"}}; user.json() }}`,
			expect: `{"address":{"city":"Rome","street":"Via Emilio Morosini"}}`,
		},
	}

	for i, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, i)
	}
}
