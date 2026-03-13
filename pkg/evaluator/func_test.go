package evaluator

import (
	"testing"

	"github.com/textwire/textwire/v3/pkg/fail"
	"github.com/textwire/textwire/v3/pkg/value"
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
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncMissingArg,
				value.STR_VAL,
				"format",
			),
		},
		// string slice
		{
			20,
			`{{ [1, 2].slice() }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncMissingArg,
				value.ARR_VAL,
				"slice",
			),
		},
		{
			30,
			`{{ [1, 2].slice("hi") }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncFirstArgInt,
				value.ARR_VAL,
				"slice",
			),
		},
		{
			40,
			`{{ [1, 2].slice({}) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncFirstArgInt,
				value.ARR_VAL,
				"slice",
			),
		},
		{
			50,
			`{{ [1, 2].slice([]) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncFirstArgInt,
				value.ARR_VAL,
				"slice",
			),
		},
		{
			60,
			`{{ [1, 2].slice(3.0) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncFirstArgInt,
				value.ARR_VAL,
				"slice",
			),
		},
		{
			70,
			`{{ [1, 2].slice(nil) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncFirstArgInt,
				value.ARR_VAL,
				"slice",
			),
		},
		{
			80,
			`{{ [1, 2].slice("hi", "hi") }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncFirstArgInt,
				value.ARR_VAL,
				"slice",
			),
		},
		{
			90,
			`{{ [1, 2].slice(0, "hi") }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncSecondArgInt,
				value.ARR_VAL,
				"slice",
			),
		},
		{
			100,
			`{{ [1, 2].slice(0, {}) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncSecondArgInt,
				value.ARR_VAL,
				"slice",
			),
		},
		{
			110,
			`{{ [1, 2].slice(0, []) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncSecondArgInt,
				value.ARR_VAL,
				"slice",
			),
		},
		{
			120,
			`{{ [1, 2].slice(0, 3.0) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncSecondArgInt,
				value.ARR_VAL,
				"slice",
			),
		},
		{
			130,
			`{{ [1, 2].slice(0, nil) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncSecondArgInt,
				value.ARR_VAL,
				"slice",
			),
		},
		// string join
		{
			140,
			`{{ [1, 2].join(1) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncFirstArgStr,
				value.ARR_VAL,
				"join",
			),
		},
		{
			150,
			`{{ [1, 2].join({}) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncFirstArgStr,
				value.ARR_VAL,
				"join",
			),
		},
		{
			160,
			`{{ [1, 2].join([]) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncFirstArgStr,
				value.ARR_VAL,
				"join",
			),
		},
		{
			170,
			`{{ [1, 2].join(3.0) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncFirstArgStr,
				value.ARR_VAL,
				"join",
			),
		},
		{
			180,
			`{{ [1, 2].join(nil) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncFirstArgStr,
				value.ARR_VAL,
				"join",
			),
		},
		// string split
		{
			190,
			`{{ "nice".split(1) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncFirstArgStr,
				value.STR_VAL,
				"split",
			),
		},
		{
			200,
			`{{ "nice".split({}) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncFirstArgStr,
				value.STR_VAL,
				"split",
			),
		},
		{
			210,
			`{{ "nice".split([]) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncFirstArgStr,
				value.STR_VAL,
				"split",
			),
		},
		{
			220,
			`{{ "nice".split(3.0) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncFirstArgStr,
				value.STR_VAL,
				"split",
			),
		},
		{
			230,
			`{{ "nice".split(nil) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncFirstArgStr,
				value.STR_VAL,
				"split",
			),
		},
		// string trim
		{
			240,
			`{{ "n".trim(1) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncFirstArgStr,
				value.STR_VAL,
				"trim",
			),
		},
		{
			250,
			`{{ "n".trim({}) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncFirstArgStr,
				value.STR_VAL,
				"trim",
			),
		},
		{
			260,
			`{{ "n".trim([]) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncFirstArgStr,
				value.STR_VAL,
				"trim",
			),
		},
		{
			270,
			`{{ "n".trim(3.0) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncFirstArgStr,
				value.STR_VAL,
				"trim",
			),
		},
		{
			280,
			`{{ "n".trim(nil) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncFirstArgStr,
				value.STR_VAL,
				"trim",
			),
		},
		// string trimRight
		{
			290,
			`{{ "n".trimRight(1) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncFirstArgStr,
				value.STR_VAL,
				"trimRight",
			),
		},
		{
			300,
			`{{ "n".trimRight({}) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncFirstArgStr,
				value.STR_VAL,
				"trimRight",
			),
		},
		{
			310,
			`{{ "n".trimRight([]) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncFirstArgStr,
				value.STR_VAL,
				"trimRight",
			),
		},
		{
			320,
			`{{ "n".trimRight(3.0) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncFirstArgStr,
				value.STR_VAL,
				"trimRight",
			),
		},
		{
			330,
			`{{ "n".trimRight(nil) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncFirstArgStr,
				value.STR_VAL,
				"trimRight",
			),
		},
		// string trimLeft
		{
			340,
			`{{ "n".trimLeft(1) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncFirstArgStr,
				value.STR_VAL,
				"trimLeft",
			),
		},
		{
			350,
			`{{ "n".trimLeft({}) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncFirstArgStr,
				value.STR_VAL,
				"trimLeft",
			),
		},
		{
			360,
			`{{ "n".trimLeft([]) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncFirstArgStr,
				value.STR_VAL,
				"trimLeft",
			),
		},
		{
			370,
			`{{ "n".trimLeft(3.0) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncFirstArgStr,
				value.STR_VAL,
				"trimLeft",
			),
		},
		{
			380,
			`{{ "n".trimLeft(nil) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncFirstArgStr,
				value.STR_VAL,
				"trimLeft",
			),
		},
		// string repeat
		{
			390,
			`{{ "n".repeat(true) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncFirstArgInt,
				value.STR_VAL,
				"repeat",
			),
		},
		{
			400,
			`{{ "n".repeat(false) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncFirstArgInt,
				value.STR_VAL,
				"repeat",
			),
		},
		{
			410,
			`{{ "n".repeat(nil) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncFirstArgInt,
				value.STR_VAL,
				"repeat",
			),
		},
		{
			420,
			`{{ "n".repeat("3") }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncFirstArgInt,
				value.STR_VAL,
				"repeat",
			),
		},
		{
			430,
			`{{ "n".repeat() }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncMissingArg,
				value.STR_VAL,
				"repeat",
			),
		},
		// string contains
		{
			440,
			`{{ "anna".contains() }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncMissingArg,
				value.STR_VAL,
				"contains",
			),
		},
		// string truncate
		{
			450,
			`{{ "anna serhii".truncate() }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncMissingArg,
				value.STR_VAL,
				"truncate",
			),
		},
		{
			460,
			`{{ "anna".truncate("hi") }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncFirstArgInt,
				value.STR_VAL,
				"truncate",
			),
		},
		{
			470,
			`{{ "anna".truncate(true) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncFirstArgInt,
				value.STR_VAL,
				"truncate",
			),
		},
		{
			480,
			`{{ "anna".truncate([]) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncFirstArgInt,
				value.STR_VAL,
				"truncate",
			),
		},
		{
			490,
			`{{ "anna".truncate({}) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncFirstArgInt,
				value.STR_VAL,
				"truncate",
			),
		},
		{
			500,
			`{{ "anna".truncate(3.3) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncFirstArgInt,
				value.STR_VAL,
				"truncate",
			),
		},
		{
			510,
			`{{ "anna".truncate(1, true) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncSecondArgStr,
				value.STR_VAL,
				"truncate",
			),
		},
		{
			520,
			`{{ "anna".truncate(2, []) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncSecondArgStr,
				value.STR_VAL,
				"truncate",
			),
		},
		{
			530,
			`{{ "anna".truncate(1, {}) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncSecondArgStr,
				value.STR_VAL,
				"truncate",
			),
		},
		{
			540,
			`{{ "anna".truncate(1, 3.3) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncSecondArgStr,
				value.STR_VAL,
				"truncate",
			),
		},
		// string decimal
		{
			550,
			`{{ "100".decimal(1) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncFirstArgStr,
				value.STR_VAL,
				"decimal",
			),
		},
		{
			560,
			`{{ "100".decimal(true) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncFirstArgStr,
				value.STR_VAL,
				"decimal",
			),
		},
		{
			570,
			`{{ "100".decimal([]) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncFirstArgStr,
				value.STR_VAL,
				"decimal",
			),
		},
		{
			580,
			`{{ "100".decimal({}) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncFirstArgStr,
				value.STR_VAL,
				"decimal",
			),
		},
		{
			590,
			`{{ "100".decimal(1.1) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncFirstArgStr,
				value.STR_VAL,
				"decimal",
			),
		},
		{
			600,
			`{{ "100".decimal("", "nice") }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncSecondArgInt,
				value.STR_VAL,
				"decimal",
			),
		},
		{
			610,
			`{{ "100".decimal("", true) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncSecondArgInt,
				value.STR_VAL,
				"decimal",
			),
		},
		{
			620,
			`{{ "100".decimal("", []) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncSecondArgInt,
				value.STR_VAL,
				"decimal",
			),
		},
		{
			630,
			`{{ "100".decimal("", {}) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncSecondArgInt,
				value.STR_VAL,
				"decimal",
			),
		},
		{
			640,
			`{{ "100".decimal("", 1.1) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncSecondArgInt,
				value.STR_VAL,
				"decimal",
			),
		},
		// integer decimal
		{
			650,
			`{{ 100.decimal(1) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncFirstArgStr,
				value.INT_VAL,
				"decimal",
			),
		},
		{
			660,
			`{{ 100.decimal(true) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncFirstArgStr,
				value.INT_VAL,
				"decimal",
			),
		},
		{
			670,
			`{{ 100.decimal([]) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncFirstArgStr,
				value.INT_VAL,
				"decimal",
			),
		},
		{
			680,
			`{{ 100.decimal({}) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncFirstArgStr,
				value.INT_VAL,
				"decimal",
			),
		},
		{
			690,
			`{{ 100.decimal(1.1) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncFirstArgStr,
				value.INT_VAL,
				"decimal",
			),
		},
		{
			700,
			`{{ 100.decimal("", "nice") }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncSecondArgInt,
				value.INT_VAL,
				"decimal",
			),
		},
		{
			710,
			`{{ 100.decimal("", true) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncSecondArgInt,
				value.INT_VAL,
				"decimal",
			),
		},
		{
			720,
			`{{ 100.decimal("", []) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncSecondArgInt,
				value.INT_VAL,
				"decimal",
			),
		},
		{
			730,
			`{{ 100.decimal("", {}) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncSecondArgInt,
				value.INT_VAL,
				"decimal",
			),
		},
		{
			740,
			`{{ 100.decimal("", 1.1) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncSecondArgInt,
				value.INT_VAL,
				"decimal",
			),
		},
		// boolean then
		{
			750,
			`{{ true.then() }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncMissingArg,
				value.BOOL_VAL,
				"then",
			),
		},
		{
			760,
			`{{ false.then() }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncMissingArg,
				value.BOOL_VAL,
				"then",
			),
		},
		// arr contains
		{
			770,
			`{{ [1, 2].contains() }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncMissingArg,
				value.ARR_VAL,
				"contains",
			),
		},
		// arr append
		{
			780,
			`{{ [1, 2].append() }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncMissingArg,
				value.ARR_VAL,
				"append",
			),
		},
		// arr prepend
		{
			790,
			`{{ [1, 2].prepend() }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncMissingArg,
				value.ARR_VAL,
				"prepend",
			),
		},
		// object get
		{
			800,
			`{{ {}.get() }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncMissingArg,
				value.OBJ_VAL,
				"get",
			),
		},
		{
			810,
			`{{ {}.get(1) }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncFirstArgStr,
				value.OBJ_VAL,
				"get",
			),
		},
		{
			820,
			`{{ {}.get("one.two", "three") }}`,
			fail.New(
				nil,
				"/path/to/file",
				fail.OriginEval,
				fail.ErrFuncMaxArgs,
				value.OBJ_VAL,
				"get",
				1,
			),
		},
	}

	for _, tc := range cases {
		evaluated, failure := testEval(tc.inp)
		if failure != nil {
			t.Fatalf("Case: %d. evaluation failed: %s", tc.id, failure)
		}

		errObj, ok := evaluated.(*value.Error)

		if !ok {
			t.Fatalf("Case: %d. Evaluation failed: %s", tc.id, errObj.String())
		}

		if evaluated.Type() != value.ERR_VAL {
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
