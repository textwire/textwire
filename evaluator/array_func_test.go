package evaluator

import (
	"testing"
)

func TestEvalArrayFunctions(t *testing.T) {
	tests := []struct {
		inp      string
		expected string
	}{
		// len
		{`{{ [].len() }}`, "0"},
		{`{{ [1, 2, 3].len() }}`, "3"},
		{`{{ [0, [2, [1, 2]]].len() }}`, "2"},
		// join
		{`{{ [1, 2, 3].join(", ") }}`, "1, 2, 3"},
		{`{{ ["one", "two", "three"].join(" ") }}`, "one two three"},
		{`{{ ["one", "two", "three"].join() }}`, "one,two,three"},
		{`{{ [].join() }}`, ""},
		// rand
		{`{{ [].rand() }}`, ""},
		{`{{ [1].rand() }}`, "1"},
		{`{{ ["some"].rand() }}`, "some"},
		{`{{ [[[4]]].rand().rand().rand() }}`, "4"},
		// reverse
		{`{{ [1, 2, 3].reverse() }}`, "3, 2, 1"},
		{`{{ ["str"].reverse() }}`, "str"},
		{`{{ [].reverse() }}`, ""},
		{`{{ ["three", "two", "one"].reverse() }}`, "one, two, three"},
		{`{{ [4, 3, [1, 2]].reverse() }}`, "1, 2, 3, 4"},
		// slice
		{`{{ [].slice(0) }}`, ""},
		{`{{ [1, 2, 3].slice(0) }}`, "1, 2, 3"},
		{`{{ [1, 2, 3].slice(-34) }}`, "1, 2, 3"}, // should change -35 to 0
		{`{{ [0, 1, 2, 3].slice(2) }}`, "2, 3"},
		{`{{ [0, 1, 2, 3].slice(5) }}`, ""},
		{`{{ [0, 1, 2, 3].slice(1, 3) }}`, "1, 2"},
		{`{{ ['one', 'two', 'three', "four"].slice(1, 2) }}`, "two"},
		{`{{ ['one', 'two', 'three', "four"].slice(0, -3) }}`, "one, two, three, four"},
		{`{{ [1, 2, 3, 4].slice(-3, -1) }}`, "1, 2, 3, 4"},
		// shuffle
		{`{{ [].shuffle() }}`, ""},
		{`{{ [1].shuffle() }}`, "1"},
		{`{{ ["test"].shuffle() }}`, "test"},
		// contains
		{`{{ [].contains(1) }}`, "0"},
		{`{{ [1, 2, 3].contains(1) }}`, "1"},
		{`{{ [1, 2, 3].contains(4) }}`, "0"},
		{`{{ [{}, [1, 2], 3].contains(1) }}`, "0"},
	}

	for _, tc := range tests {
		evaluationExpected(t, tc.inp, tc.expected)
	}
}
