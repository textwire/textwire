package evaluator

import "testing"

func TestEvalArrayFunctions(t *testing.T) {
	tests := []struct {
		inp      string
		expected string
	}{
		// len func
		{`{{ [].len() }}`, "0"},
		{`{{ [1, 2, 3].len() }}`, "3"},
		{`{{ [0, [2, [1, 2]]].len() }}`, "2"},
		// join func
		{`{{ [1, 2, 3].join(", ") }}`, "1, 2, 3"},
		{`{{ ["one", "two", "three"].join(" ") }}`, "one two three"},
		{`{{ ["one", "two", "three"].join() }}`, "one,two,three"},
		{`{{ [].join() }}`, ""},
		// rand func
		{`{{ [].rand() }}`, ""},
		{`{{ [1].rand() }}`, "1"},
		{`{{ ["some"].rand() }}`, "some"},
		{`{{ [[[4]]].rand().rand().rand() }}`, "4"},
		// reverse func
		{`{{ [1, 2, 3].reverse() }}`, "3, 2, 1"},
		{`{{ ["str"].reverse() }}`, "str"},
		{`{{ [].reverse() }}`, ""},
		{`{{ ["three", "two", "one"].reverse() }}`, "one, two, three"},
		{`{{ [4, 3, [1, 2]].reverse() }}`, "1, 2, 3, 4"},
	}

	for _, tc := range tests {
		evaluationExpected(t, tc.inp, tc.expected)
	}
}
