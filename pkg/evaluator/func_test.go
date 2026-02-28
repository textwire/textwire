package evaluator

import (
	"testing"

	"github.com/textwire/textwire/v3/pkg/fail"
	"github.com/textwire/textwire/v3/pkg/object"
)

func TestFunctionGivesError(t *testing.T) {
	cases := []struct {
		id  int
		inp string
		err *fail.Error
	}{
		// string format
		{
			10,
			`{{ "He has %s apples".format() }}`,
			fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrFuncMissingArg,
				object.STR_OBJ,
				"format",
			),
		},
		// string slice
		{
			20,
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
			30,
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
			40,
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
			50,
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
			60,
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
			70,
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
			80,
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
			90,
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
			100,
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
			110,
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
			120,
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
			130,
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
		// string join
		{
			140,
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
			150,
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
			160,
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
			170,
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
			180,
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
		// string split
		{
			190,
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
			200,
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
			210,
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
			220,
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
			230,
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
		// string trim
		{
			240,
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
			250,
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
			260,
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
			270,
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
			280,
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
		// string trimRight
		{
			290,
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
			300,
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
			310,
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
			320,
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
			330,
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
		// string trimLeft
		{
			340,
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
			350,
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
			360,
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
			370,
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
			380,
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
		// string repeat
		{
			390,
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
			400,
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
			410,
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
			420,
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
			430,
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
		// string contains
		{
			440,
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
		// string truncate
		{
			450,
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
			460,
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
			470,
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
			480,
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
			490,
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
			500,
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
			510,
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
			520,
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
			530,
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
			540,
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
		// string decimal
		{
			550,
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
			560,
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
			570,
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
			580,
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
			590,
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
			600,
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
			610,
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
			620,
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
			630,
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
			640,
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
		// integer decimal
		{
			650,
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
			660,
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
			670,
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
			680,
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
			690,
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
			700,
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
			710,
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
			720,
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
			730,
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
			740,
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
		// boolean then
		{
			750,
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
			760,
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
		// array contains
		{
			770,
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
		// array append
		{
			780,
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
		// array prepend
		{
			790,
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
			t.Fatalf("Case: %d. Evaluation failed: %s", tc.id, errObj.String())
		}

		if evaluated.Type() != object.ERR_OBJ {
			t.Fatalf("Case: %d. Expect object.ERR_OBJ, got=%T", tc.id, evaluated)
		}

		if errObj.String() != tc.err.String() {
			t.Fatalf(
				"Case: %d. Expect error message=%q, got=%q",
				tc.id,
				tc.err.String(),
				errObj.String(),
			)
		}
	}
}
