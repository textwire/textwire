package evaluator

import (
	"testing"

	"github.com/textwire/textwire/v2/fail"
	"github.com/textwire/textwire/v2/object"
)

func TestFunctionGivesError(t *testing.T) {
	tests := []struct {
		inp string
		err *fail.Error
	}{
		// slice
		{`{{ [1, 2].slice() }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncRequiresOneArg, "slice", "array")},
		{`{{ [1, 2].slice("hi") }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgInt, "slice", "array")},
		{`{{ [1, 2].slice({}) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgInt, "slice", "array")},
		{`{{ [1, 2].slice([]) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgInt, "slice", "array")},
		{`{{ [1, 2].slice(3.0) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgInt, "slice", "array")},
		{`{{ [1, 2].slice(nil) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgInt, "slice", "array")},
		{`{{ [1, 2].slice("hi", "hi") }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgInt, "slice", "array")},
		{`{{ [1, 2].slice(0, "hi") }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncSecondArgInt, "slice", "array")},
		{`{{ [1, 2].slice(0, {}) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncSecondArgInt, "slice", "array")},
		{`{{ [1, 2].slice(0, []) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncSecondArgInt, "slice", "array")},
		{`{{ [1, 2].slice(0, 3.0) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncSecondArgInt, "slice", "array")},
		{`{{ [1, 2].slice(0, nil) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncSecondArgInt, "slice", "array")},
		// join
		{`{{ [1, 2].join(1) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgStr, "join", "array")},
		{`{{ [1, 2].join({}) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgStr, "join", "array")},
		{`{{ [1, 2].join([]) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgStr, "join", "array")},
		{`{{ [1, 2].join(3.0) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgStr, "join", "array")},
		{`{{ [1, 2].join(nil) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgStr, "join", "array")},
		// split
		{`{{ "nice".split(1) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgStr, "split", "string")},
		{`{{ "nice".split({}) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgStr, "split", "string")},
		{`{{ "nice".split([]) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgStr, "split", "string")},
		{`{{ "nice".split(3.0) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgStr, "split", "string")},
		{`{{ "nice".split(nil) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgStr, "split", "string")},
		// trim
		{`{{ " nice".trim(1) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgStr, "trim", "string")},
		{`{{ " nice".trim({}) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgStr, "trim", "string")},
		{`{{ " nice".trim([]) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgStr, "trim", "string")},
		{`{{ " nice".trim(3.0) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgStr, "trim", "string")},
		{`{{ " nice".trim(nil) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgStr, "trim", "string")},
		// contains
		{`{{ "anna".contains() }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncRequiresOneArg, "contains", "string")},
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

		if errObj.String() != tc.err.String() {
			t.Fatalf("expected error message=%q, got=%q", tc.err.String(), errObj.String())
		}
	}
}
