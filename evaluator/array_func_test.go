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
		{`{{ [1, 2, 3].contains(2) }}`, "1"},
		{`{{ [1, 2, 3].contains(3) }}`, "1"},
		{`{{ [1, 2, 3].contains(4) }}`, "0"},
		{`{{ [{}, [1, 2], 3].contains(1) }}`, "0"},
		{`{{ [{}, [1, 2], 3].contains({}) }}`, "1"},
		{`{{ [{nice: 3, cool: 'anna'}, [1, 2], 3].contains({nice:3, cool: 'anna'}) }}`, "1"},
		{`{{ [{cool: 'anna', nice: 3}, [1, 2], 3].contains({nice:3, cool: 'anna'}) }}`, "1"},
		{`{{ [[], 3].contains([]) }}`, "1"},
		{`{{ [[1, 2], 3].contains([1, 2]) }}`, "1"},
		{`{{ [[2, 1], 3].contains([1, 2]) }}`, "0"},
		{`{{ [1, 2].contains({}) }}`, "0"},
		{`{{ [{}, 21].contains({age: 21}) }}`, "0"},
		{`{{ [[], [1], [2]].contains([2]) }}`, "1"},
		// append
		{`{{ [].append(1) }}`, "1"},
		{`{{ [1, 2, 3].append(4) }}`, "1, 2, 3, 4"},
		{`{{ ["one", "two"].append("three") }}`, "one, two, three"},
		{`{{ [1, 2].append(3, 4) }}`, "1, 2, 3, 4"},
		{`{{ [1, 2].append([3, 4]) }}`, "1, 2, 3, 4"},
		{`{{ [1, 2].append([3, 4]).len() }}`, "3"},
		// prepend
		{`{{ [].prepend(1) }}`, "1"},
		{`{{ [1, 2, 3].prepend(4) }}`, "4, 1, 2, 3"},
		{`{{ ["one", "two"].prepend("three") }}`, "three, one, two"},
		{`{{ [1, 2].prepend(3, 4) }}`, "3, 4, 1, 2"},
		{`{{ [1, 2].prepend([3, 4]) }}`, "3, 4, 1, 2"},
		{`{{ [1, 2].prepend([3, 4]).len() }}`, "3"},
	}

	for _, tc := range tests {
		evaluationExpected(t, tc.inp, tc.expected)
	}
}
