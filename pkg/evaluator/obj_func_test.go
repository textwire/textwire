package evaluator

import (
	"testing"
)

func TestEvalObjectFunctions(t *testing.T) {
	cases := []struct {
		id     int
		inp    string
		expect string
	}{
		// json
		{10, `{{ {}.json() }}`, "{}"},
		{20, `{{ {one: {two: {}}}.json() }}`, `{"one":{"two":{}}}`},
		{30, `{{ {name: "Chiori", game: "Genshin Impact"}.json() }}`, `{"game":"Genshin Impact","name":"Chiori"}`},
		{40, `{{ user = {address: {street: "Via Emilio Morosini", city: "Rome"}}; user.json() }}`, `{"address":{"city":"Rome","street":"Via Emilio Morosini"}}`},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, tc.id)
	}
}
