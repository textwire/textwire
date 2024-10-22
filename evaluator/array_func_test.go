package evaluator

import (
	"strings"
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
	}

	for _, tc := range tests {
		evaluationExpected(t, tc.inp, tc.expected)
	}
}

func TestArrayShuffle(t *testing.T) {
	input := `{{ [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15].shuffle() }}`
	initial := "0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15"

	shuffled := testEval(input).String()

	// Ensure the shuffled array has the same length as the original
	if len(shuffled) != len(initial) {
		t.Errorf("expected shuffled array to have length %d, got %d", len(initial), len(shuffled))
	}

	// Check that the shuffle changed the order
	if shuffled == initial {
		t.Error("expected shuffled array to be different from the initial array")
	}

	// Check that all numbers from the initial array are present in the shuffled array
	initialElements := strings.Split(initial, ", ")
	shuffledElements := strings.Split(shuffled, ", ")

	// Create a map to track occurrences of elements in the shuffled array
	elementMap := make(map[string]bool)

	for _, elem := range shuffledElements {
		elementMap[elem] = true
	}

	// Ensure all elements from the initial array are in the shuffled array
	for _, elem := range initialElements {
		if !elementMap[elem] {
			t.Errorf("expected element %s to be present in shuffled array", elem)
		}
	}
}
