package evaluator

import (
	"testing"

	"github.com/textwire/textwire/v3/fail"
	"github.com/textwire/textwire/v3/object"
)

func TestFunctionGivesError(t *testing.T) {
	cases := []struct {
		inp string
		err *fail.Error
	}{
		// slice
		{
			`{{ [1, 2].slice() }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncMissingArg,
				object.ARR_OBJ,
				"slice",
			),
		},
		{
			`{{ [1, 2].slice("hi") }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncFirstArgInt,
				object.ARR_OBJ,
				"slice",
			),
		},
		{
			`{{ [1, 2].slice({}) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncFirstArgInt,
				object.ARR_OBJ,
				"slice",
			),
		},
		{
			`{{ [1, 2].slice([]) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncFirstArgInt,
				object.ARR_OBJ,
				"slice",
			),
		},
		{
			`{{ [1, 2].slice(3.0) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncFirstArgInt,
				object.ARR_OBJ,
				"slice",
			),
		},
		{
			`{{ [1, 2].slice(nil) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncFirstArgInt,
				object.ARR_OBJ,
				"slice",
			),
		},
		{
			`{{ [1, 2].slice("hi", "hi") }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncFirstArgInt,
				object.ARR_OBJ,
				"slice",
			),
		},
		{
			`{{ [1, 2].slice(0, "hi") }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncSecondArgInt,
				object.ARR_OBJ,
				"slice",
			),
		},
		{
			`{{ [1, 2].slice(0, {}) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncSecondArgInt,
				object.ARR_OBJ,
				"slice",
			),
		},
		{
			`{{ [1, 2].slice(0, []) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncSecondArgInt,
				object.ARR_OBJ,
				"slice",
			),
		},
		{
			`{{ [1, 2].slice(0, 3.0) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncSecondArgInt,
				object.ARR_OBJ,
				"slice",
			),
		},
		{
			`{{ [1, 2].slice(0, nil) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncSecondArgInt,
				object.ARR_OBJ,
				"slice",
			),
		},
		// join
		{
			`{{ [1, 2].join(1) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncFirstArgStr,
				object.ARR_OBJ,
				"join",
			),
		},
		{
			`{{ [1, 2].join({}) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncFirstArgStr,
				object.ARR_OBJ,
				"join",
			),
		},
		{
			`{{ [1, 2].join([]) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncFirstArgStr,
				object.ARR_OBJ,
				"join",
			),
		},
		{
			`{{ [1, 2].join(3.0) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncFirstArgStr,
				object.ARR_OBJ,
				"join",
			),
		},
		{
			`{{ [1, 2].join(nil) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncFirstArgStr,
				object.ARR_OBJ,
				"join",
			),
		},
		// split
		{
			`{{ "nice".split(1) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncFirstArgStr,
				object.STR_OBJ,
				"split",
			),
		},
		{
			`{{ "nice".split({}) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncFirstArgStr,
				object.STR_OBJ,
				"split",
			),
		},
		{
			`{{ "nice".split([]) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncFirstArgStr,
				object.STR_OBJ,
				"split",
			),
		},
		{
			`{{ "nice".split(3.0) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncFirstArgStr,
				object.STR_OBJ,
				"split",
			),
		},
		{
			`{{ "nice".split(nil) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncFirstArgStr,
				object.STR_OBJ,
				"split",
			),
		},
		// trim
		{
			`{{ "n".trim(1) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncFirstArgStr,
				object.STR_OBJ,
				"trim",
			),
		},
		{
			`{{ "n".trim({}) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncFirstArgStr,
				object.STR_OBJ,
				"trim",
			),
		},
		{
			`{{ "n".trim([]) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncFirstArgStr,
				object.STR_OBJ,
				"trim",
			),
		},
		{
			`{{ "n".trim(3.0) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncFirstArgStr,
				object.STR_OBJ,
				"trim",
			),
		},
		{
			`{{ "n".trim(nil) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncFirstArgStr,
				object.STR_OBJ,
				"trim",
			),
		},
		// trimRight
		{
			`{{ "n".trimRight(1) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncFirstArgStr,
				object.STR_OBJ,
				"trimRight",
			),
		},
		{
			`{{ "n".trimRight({}) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncFirstArgStr,
				object.STR_OBJ,
				"trimRight",
			),
		},
		{
			`{{ "n".trimRight([]) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncFirstArgStr,
				object.STR_OBJ,
				"trimRight",
			),
		},
		{
			`{{ "n".trimRight(3.0) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncFirstArgStr,
				object.STR_OBJ,
				"trimRight",
			),
		},
		{
			`{{ "n".trimRight(nil) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncFirstArgStr,
				object.STR_OBJ,
				"trimRight",
			),
		},
		// trimLeft
		{
			`{{ "n".trimLeft(1) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncFirstArgStr,
				object.STR_OBJ,
				"trimLeft",
			),
		},
		{
			`{{ "n".trimLeft({}) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncFirstArgStr,
				object.STR_OBJ,
				"trimLeft",
			),
		},
		{
			`{{ "n".trimLeft([]) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncFirstArgStr,
				object.STR_OBJ,
				"trimLeft",
			),
		},
		{
			`{{ "n".trimLeft(3.0) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncFirstArgStr,
				object.STR_OBJ,
				"trimLeft",
			),
		},
		{
			`{{ "n".trimLeft(nil) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncFirstArgStr,
				object.STR_OBJ,
				"trimLeft",
			),
		},
		// repeat
		{
			`{{ "n".repeat(true) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncFirstArgInt,
				object.STR_OBJ,
				"repeat",
			),
		},
		{
			`{{ "n".repeat(false) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncFirstArgInt,
				object.STR_OBJ,
				"repeat",
			),
		},
		{
			`{{ "n".repeat(nil) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncFirstArgInt,
				object.STR_OBJ,
				"repeat",
			),
		},
		{
			`{{ "n".repeat("3") }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncFirstArgInt,
				object.STR_OBJ,
				"repeat",
			),
		},
		{
			`{{ "n".repeat() }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncMissingArg,
				object.STR_OBJ,
				"repeat",
			),
		},
		// contains
		{
			`{{ "anna".contains() }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncMissingArg,
				object.STR_OBJ,
				"contains",
			),
		},
		// truncate
		{
			`{{ "anna serhii".truncate() }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncMissingArg,
				object.STR_OBJ,
				"truncate",
			),
		},
		{
			`{{ "anna".truncate("hi") }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncFirstArgInt,
				object.STR_OBJ,
				"truncate",
			),
		},
		{
			`{{ "anna".truncate(true) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncFirstArgInt,
				object.STR_OBJ,
				"truncate",
			),
		},
		{
			`{{ "anna".truncate([]) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncFirstArgInt,
				object.STR_OBJ,
				"truncate",
			),
		},
		{
			`{{ "anna".truncate({}) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncFirstArgInt,
				object.STR_OBJ,
				"truncate",
			),
		},
		{
			`{{ "anna".truncate(3.3) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncFirstArgInt,
				object.STR_OBJ,
				"truncate",
			),
		},
		{
			`{{ "anna".truncate(1, true) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncSecondArgStr,
				object.STR_OBJ,
				"truncate",
			),
		},
		{
			`{{ "anna".truncate(2, []) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncSecondArgStr,
				object.STR_OBJ,
				"truncate",
			),
		},
		{
			`{{ "anna".truncate(1, {}) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncSecondArgStr,
				object.STR_OBJ,
				"truncate",
			),
		},
		{
			`{{ "anna".truncate(1, 3.3) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncSecondArgStr,
				object.STR_OBJ,
				"truncate",
			),
		},
		// decimal (STRING)
		{
			`{{ "100".decimal(1) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncFirstArgStr,
				object.STR_OBJ,
				"decimal",
			),
		},
		{
			`{{ "100".decimal(true) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncFirstArgStr,
				object.STR_OBJ,
				"decimal",
			),
		},
		{
			`{{ "100".decimal([]) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncFirstArgStr,
				object.STR_OBJ,
				"decimal",
			),
		},
		{
			`{{ "100".decimal({}) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncFirstArgStr,
				object.STR_OBJ,
				"decimal",
			),
		},
		{
			`{{ "100".decimal(1.1) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncFirstArgStr,
				object.STR_OBJ,
				"decimal",
			),
		},
		{
			`{{ "100".decimal("", "nice") }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncSecondArgInt,
				object.STR_OBJ,
				"decimal",
			),
		},
		{
			`{{ "100".decimal("", true) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncSecondArgInt,
				object.STR_OBJ,
				"decimal",
			),
		},
		{
			`{{ "100".decimal("", []) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncSecondArgInt,
				object.STR_OBJ,
				"decimal",
			),
		},
		{
			`{{ "100".decimal("", {}) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncSecondArgInt,
				object.STR_OBJ,
				"decimal",
			),
		},
		{
			`{{ "100".decimal("", 1.1) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncSecondArgInt,
				object.STR_OBJ,
				"decimal",
			),
		},
		// decimal (INTEGER)
		{
			`{{ 100.decimal(1) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncFirstArgStr,
				object.INT_OBJ,
				"decimal",
			),
		},
		{
			`{{ 100.decimal(true) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncFirstArgStr,
				object.INT_OBJ,
				"decimal",
			),
		},
		{
			`{{ 100.decimal([]) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncFirstArgStr,
				object.INT_OBJ,
				"decimal",
			),
		},
		{
			`{{ 100.decimal({}) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncFirstArgStr,
				object.INT_OBJ,
				"decimal",
			),
		},
		{
			`{{ 100.decimal(1.1) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncFirstArgStr,
				object.INT_OBJ,
				"decimal",
			),
		},
		{
			`{{ 100.decimal("", "nice") }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncSecondArgInt,
				object.INT_OBJ,
				"decimal",
			),
		},
		{
			`{{ 100.decimal("", true) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncSecondArgInt,
				object.INT_OBJ,
				"decimal",
			),
		},
		{
			`{{ 100.decimal("", []) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncSecondArgInt,
				object.INT_OBJ,
				"decimal",
			),
		},
		{
			`{{ 100.decimal("", {}) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncSecondArgInt,
				object.INT_OBJ,
				"decimal",
			),
		},
		{
			`{{ 100.decimal("", 1.1) }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncSecondArgInt,
				object.INT_OBJ,
				"decimal",
			),
		},
		// then
		{
			`{{ true.then() }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncMissingArg,
				object.BOOL_OBJ,
				"then",
			),
		},
		{
			`{{ false.then() }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncMissingArg,
				object.BOOL_OBJ,
				"then",
			),
		},
		// contains
		{
			`{{ [1, 2].contains() }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncMissingArg,
				object.ARR_OBJ,
				"contains",
			),
		},
		// append
		{
			`{{ [1, 2].append() }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncMissingArg,
				object.ARR_OBJ,
				"append",
			),
		},
		// prepend
		{
			`{{ [1, 2].prepend() }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncMissingArg,
				object.ARR_OBJ,
				"prepend",
			),
		},
	}

	for _, tc := range cases {
		evaluated := testEval(tc.inp)
		errObj, ok := evaluated.(*object.Error)

		if !ok {
			t.Fatalf("evaluation failed: %s", errObj.String())
		}

		if evaluated.Type() != object.ERR_OBJ {
			t.Fatalf("expect object.ERR_OBJ, got=%T", evaluated)
		}

		if errObj.String() != tc.err.String() {
			t.Fatalf("expect error message=%q, got=%q", tc.err.String(), errObj.String())
		}
	}
}
