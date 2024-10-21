package evaluator

import (
	"testing"

	"github.com/textwire/textwire/v2/fail"
	"github.com/textwire/textwire/v2/object"
)

func TestFunctionGivesError(t *testing.T) {
	tests := []struct {
		inp         string
		expectedErr string
		funcName    string
	}{
		// slice function
		{`{{ [1, 2].slice() }}`, fail.ErrFuncRequiresOneArg, "slice"},
		{`{{ [1, 2].slice("hi") }}`, fail.ErrFuncFirstArgInt, "slice"},
		{`{{ [1, 2].slice({}) }}`, fail.ErrFuncFirstArgInt, "slice"},
		{`{{ [1, 2].slice([]) }}`, fail.ErrFuncFirstArgInt, "slice"},
		{`{{ [1, 2].slice(3.0) }}`, fail.ErrFuncFirstArgInt, "slice"},
		{`{{ [1, 2].slice(nil) }}`, fail.ErrFuncFirstArgInt, "slice"},
		{`{{ [1, 2].slice("hi", "hi") }}`, fail.ErrFuncFirstArgInt, "slice"},
		{`{{ [1, 2].slice(0, "hi") }}`, fail.ErrFuncSecondArgInt, "slice"},
		{`{{ [1, 2].slice(0, {}) }}`, fail.ErrFuncSecondArgInt, "slice"},
		{`{{ [1, 2].slice(0, []) }}`, fail.ErrFuncSecondArgInt, "slice"},
		{`{{ [1, 2].slice(0, 3.0) }}`, fail.ErrFuncSecondArgInt, "slice"},
		{`{{ [1, 2].slice(0, nil) }}`, fail.ErrFuncSecondArgInt, "slice"},
		// join function
		{`{{ [1, 2].join(1) }}`, fail.ErrFuncFirstArgStr, "join"},
		{`{{ [1, 2].join({}) }}`, fail.ErrFuncFirstArgStr, "join"},
		{`{{ [1, 2].join([]) }}`, fail.ErrFuncFirstArgStr, "join"},
		{`{{ [1, 2].join(3.0) }}`, fail.ErrFuncFirstArgStr, "join"},
		{`{{ [1, 2].join(nil) }}`, fail.ErrFuncFirstArgStr, "join"},
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
