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
		{
			30,
			`{{ {name: "Chiori", game: "Genshin Impact"}.json() }}`,
			`{"game":"Genshin Impact","name":"Chiori"}`,
		},
		{
			40,
			`{{ user = {address: {street: "Via Emilio Morosini", city: "Rome"}}; user.json() }}`,
			`{"address":{"city":"Rome","street":"Via Emilio Morosini"}}`,
		},
		{50, `{{ {a: {b: {c: {d: 1}}}}.json() }}`, `{"a":{"b":{"c":{"d":1}}}}`},
		{
			60,
			`{{ {nums: [1, 2, 3], strs: ['a', 'b']}.json() }}`,
			`{"nums":[1,2,3],"strs":["a","b"]}`,
		},
		{
			70,
			`{{ {quote: 'He said Hello', newline: 'A B'}.json() }}`,
			`{"newline":"A B","quote":"He said Hello"}`,
		},
		{
			80,
			`{{ {active: true, count: nil, rate: 3.14}.json() }}`,
			`{"active":true,"count":null,"rate":3.14}`,
		},
		{90, `{{ {z: 1, a: 2, m: 3}.json() }}`, `{"a":2,"m":3,"z":1}`},
		{
			100,
			`{{ {user: {name: 'John', age: 30, hobbies: ['coding', 'gaming']}, active: true}.json() }}`,
			`{"active":true,"user":{"age":30,"hobbies":["coding","gaming"],"name":"John"}}`,
		},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, tc.id)
	}
}
