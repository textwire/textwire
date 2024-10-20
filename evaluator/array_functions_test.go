package evaluator

import "testing"

func TestEvalArrayFunctions(t *testing.T) {
	tests := []struct {
		inp      string
		expected string
	}{
		{`{{ [].len() }}`, "0"},
		{`{{ [1, 2, 3].len() }}`, "3"},
		{`{{ [0, [2, [1, 2]]].len() }}`, "2"},
		{`{{ [1, 2, 3].join(", ") }}`, "1, 2, 3"},
		{`{{ ["one", "two", "three"].join(" ") }}`, "one two three"},
		{`{{ ["one", "two", "three"].join() }}`, "one,two,three"},
		{`{{ [].join() }}`, ""},
		{`{{ [].rand() }}`, ""},
		{`{{ [1].rand() }}`, "1"},
		{`{{ ["some"].rand() }}`, "some"},
		{`{{ [[[4]]].rand().rand().rand() }}`, "4"},
	}

	for _, tc := range tests {
		evaluationExpected(t, tc.inp, tc.expected)
	}
}
