package evaluator

import (
	"testing"

	"github.com/textwire/textwire/v2/fail"
	"github.com/textwire/textwire/v2/object"
)

func TestFunctionGivesError(t *testing.T) {
	cases := []struct {
		inp string
		err *fail.Error
	}{
		// slice
		{`{{ [1, 2].slice() }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncRequiresOneArg, "slice", object.ARR_OBJ)},
		{`{{ [1, 2].slice("hi") }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgInt, "slice", object.ARR_OBJ)},
		{`{{ [1, 2].slice({}) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgInt, "slice", object.ARR_OBJ)},
		{`{{ [1, 2].slice([]) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgInt, "slice", object.ARR_OBJ)},
		{`{{ [1, 2].slice(3.0) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgInt, "slice", object.ARR_OBJ)},
		{`{{ [1, 2].slice(nil) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgInt, "slice", object.ARR_OBJ)},
		{`{{ [1, 2].slice("hi", "hi") }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgInt, "slice", object.ARR_OBJ)},
		{`{{ [1, 2].slice(0, "hi") }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncSecondArgInt, "slice", object.ARR_OBJ)},
		{`{{ [1, 2].slice(0, {}) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncSecondArgInt, "slice", object.ARR_OBJ)},
		{`{{ [1, 2].slice(0, []) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncSecondArgInt, "slice", object.ARR_OBJ)},
		{`{{ [1, 2].slice(0, 3.0) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncSecondArgInt, "slice", object.ARR_OBJ)},
		{`{{ [1, 2].slice(0, nil) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncSecondArgInt, "slice", object.ARR_OBJ)},
		// join
		{`{{ [1, 2].join(1) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgStr, "join", object.ARR_OBJ)},
		{`{{ [1, 2].join({}) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgStr, "join", object.ARR_OBJ)},
		{`{{ [1, 2].join([]) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgStr, "join", object.ARR_OBJ)},
		{`{{ [1, 2].join(3.0) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgStr, "join", object.ARR_OBJ)},
		{`{{ [1, 2].join(nil) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgStr, "join", object.ARR_OBJ)},
		// split
		{`{{ "nice".split(1) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgStr, "split", object.STR_OBJ)},
		{`{{ "nice".split({}) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgStr, "split", object.STR_OBJ)},
		{`{{ "nice".split([]) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgStr, "split", object.STR_OBJ)},
		{`{{ "nice".split(3.0) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgStr, "split", object.STR_OBJ)},
		{`{{ "nice".split(nil) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgStr, "split", object.STR_OBJ)},
		// trim
		{`{{ "n".trim(1) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgStr, "trim", object.STR_OBJ)},
		{`{{ "n".trim({}) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgStr, "trim", object.STR_OBJ)},
		{`{{ "n".trim([]) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgStr, "trim", object.STR_OBJ)},
		{`{{ "n".trim(3.0) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgStr, "trim", object.STR_OBJ)},
		{`{{ "n".trim(nil) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgStr, "trim", object.STR_OBJ)},
		// trimRight
		{`{{ "n".trimRight(1) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgStr, "trimRight", object.STR_OBJ)},
		{`{{ "n".trimRight({}) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgStr, "trimRight", object.STR_OBJ)},
		{`{{ "n".trimRight([]) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgStr, "trimRight", object.STR_OBJ)},
		{`{{ "n".trimRight(3.0) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgStr, "trimRight", object.STR_OBJ)},
		{`{{ "n".trimRight(nil) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgStr, "trimRight", object.STR_OBJ)},
		// trimLeft
		{`{{ "n".trimLeft(1) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgStr, "trimLeft", object.STR_OBJ)},
		{`{{ "n".trimLeft({}) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgStr, "trimLeft", object.STR_OBJ)},
		{`{{ "n".trimLeft([]) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgStr, "trimLeft", object.STR_OBJ)},
		{`{{ "n".trimLeft(3.0) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgStr, "trimLeft", object.STR_OBJ)},
		{`{{ "n".trimLeft(nil) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgStr, "trimLeft", object.STR_OBJ)},
		// repeat
		{`{{ "n".repeat(true) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgInt, "repeat", object.STR_OBJ)},
		{`{{ "n".repeat(false) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgInt, "repeat", object.STR_OBJ)},
		{`{{ "n".repeat(nil) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgInt, "repeat", object.STR_OBJ)},
		{`{{ "n".repeat("3") }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgInt, "repeat", object.STR_OBJ)},
		{`{{ "n".repeat() }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncRequiresOneArg, "repeat", object.STR_OBJ)},
		// contains
		{`{{ "anna".contains() }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncRequiresOneArg, "contains", object.STR_OBJ)},
		// truncate
		{`{{ "anna serhii".truncate() }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncRequiresOneArg, "truncate", object.STR_OBJ)},
		{`{{ "anna".truncate("hi") }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgInt, "truncate", object.STR_OBJ)},
		{`{{ "anna".truncate(true) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgInt, "truncate", object.STR_OBJ)},
		{`{{ "anna".truncate([]) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgInt, "truncate", object.STR_OBJ)},
		{`{{ "anna".truncate({}) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgInt, "truncate", object.STR_OBJ)},
		{`{{ "anna".truncate(3.3) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgInt, "truncate", object.STR_OBJ)},
		{`{{ "anna".truncate(1, true) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncSecondArgStr, "truncate", object.STR_OBJ)},
		{`{{ "anna".truncate(2, []) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncSecondArgStr, "truncate", object.STR_OBJ)},
		{`{{ "anna".truncate(1, {}) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncSecondArgStr, "truncate", object.STR_OBJ)},
		{`{{ "anna".truncate(1, 3.3) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncSecondArgStr, "truncate", object.STR_OBJ)},
		// decimal (STRING)
		{`{{ "100".decimal(1) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgStr, "decimal", object.STR_OBJ)},
		{`{{ "100".decimal(true) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgStr, "decimal", object.STR_OBJ)},
		{`{{ "100".decimal([]) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgStr, "decimal", object.STR_OBJ)},
		{`{{ "100".decimal({}) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgStr, "decimal", object.STR_OBJ)},
		{`{{ "100".decimal(1.1) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgStr, "decimal", object.STR_OBJ)},
		{`{{ "100".decimal("", "nice") }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncSecondArgInt, "decimal", object.STR_OBJ)},
		{`{{ "100".decimal("", true) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncSecondArgInt, "decimal", object.STR_OBJ)},
		{`{{ "100".decimal("", []) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncSecondArgInt, "decimal", object.STR_OBJ)},
		{`{{ "100".decimal("", {}) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncSecondArgInt, "decimal", object.STR_OBJ)},
		{`{{ "100".decimal("", 1.1) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncSecondArgInt, "decimal", object.STR_OBJ)},
		// decimal (INTEGER)
		{`{{ 100.decimal(1) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgStr, "decimal", object.INT_OBJ)},
		{`{{ 100.decimal(true) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgStr, "decimal", object.INT_OBJ)},
		{`{{ 100.decimal([]) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgStr, "decimal", object.INT_OBJ)},
		{`{{ 100.decimal({}) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgStr, "decimal", object.INT_OBJ)},
		{`{{ 100.decimal(1.1) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncFirstArgStr, "decimal", object.INT_OBJ)},
		{`{{ 100.decimal("", "nice") }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncSecondArgInt, "decimal", object.INT_OBJ)},
		{`{{ 100.decimal("", true) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncSecondArgInt, "decimal", object.INT_OBJ)},
		{`{{ 100.decimal("", []) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncSecondArgInt, "decimal", object.INT_OBJ)},
		{`{{ 100.decimal("", {}) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncSecondArgInt, "decimal", object.INT_OBJ)},
		{`{{ 100.decimal("", 1.1) }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncSecondArgInt, "decimal", object.INT_OBJ)},
		// then
		{`{{ true.then() }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncRequiresOneArg, "then", object.BOOL_OBJ)},
		{`{{ false.then() }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncRequiresOneArg, "then", object.BOOL_OBJ)},
		// contains
		{`{{ [1, 2].contains() }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncRequiresOneArg, "contains", object.ARR_OBJ)},
		// append
		{`{{ [1, 2].append() }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncRequiresOneArg, "append", object.ARR_OBJ)},
		// prepend
		{`{{ [1, 2].prepend() }}`, fail.New(1, "/path/to/file", "evaluator", fail.ErrFuncRequiresOneArg, "prepend", object.ARR_OBJ)},
	}

	for _, tc := range cases {
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
