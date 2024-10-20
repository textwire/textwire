package evaluator

import (
	"testing"

	"github.com/textwire/textwire/v2/fail"
	"github.com/textwire/textwire/v2/object"
)

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
		// slice func
		{`{{ [].slice(0) }}`, ""},
		{`{{ [1, 2, 3].slice(0) }}`, "1, 2, 3"},
		{`{{ [1, 2, 3].slice(-34) }}`, "1, 2, 3"}, // should change -35 to 0
		{`{{ [0, 1, 2, 3].slice(2) }}`, "2, 3"},
		{`{{ [0, 1, 2, 3].slice(5) }}`, ""},
	}

	for _, tc := range tests {
		evaluationExpected(t, tc.inp, tc.expected)
	}
}

func TestFunctionGivesError(t *testing.T) {
	tests := []struct {
		inp         string
		expectedErr string
		funcName    string
	}{
		{`{{ [1, 2].slice() }}`, fail.ErrFuncRequiresOneArg, "slice"},
	}

	for _, tc := range tests {
		evaluated := testEval(tc.inp)
		errObj, ok := evaluated.(*object.Error)

		if !ok {
			t.Fatalf("evaluation failed: %s", errObj.String())
		}

		if evaluated.Type() != object.ERR_OBJ {
			t.Fatalf("expected object.ERR_OBJ, got=%T", evaluated)
		}

		failErr := fail.New(1, "/path/to/file", "evaluator", tc.expectedErr, tc.funcName)

		if errObj.String() != failErr.String() {
			t.Fatalf("expected error message=%q, got=%q", failErr.String(), errObj.String())
		}
	}
}
