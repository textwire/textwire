package textwire

import (
	"testing"

	"codeberg.org/textwire/textwire/v5/pkg/fail"
)

func TestCustomFunc_UnsupportedStringResult(t *testing.T) {
	returnsFunc := func(string, ...any) any { return unsupportedFunc }
	returnsChan := func(string, ...any) any { return unsupportedChan }

	assertNoFailure(t, RegisterStrFunc("_returnsFunc", returnsFunc))
	assertNoFailure(t, RegisterStrFunc("_returnsChan", returnsChan))

	_, failure := EvaluateString(`{{ "x"._returnsFunc() }}`, nil)
	assertUnsupportedResult(t, failure, unsupportedFunc)

	_, failure = EvaluateString(`{{ "x"._returnsChan() }}`, nil)
	assertUnsupportedResult(t, failure, unsupportedChan)
}

func TestCustomFunc_UnsupportedArrayResult(t *testing.T) {
	returnsFunc := func([]any, ...any) any { return unsupportedFunc }
	returnsChan := func([]any, ...any) any { return unsupportedChan }

	assertNoFailure(t, RegisterArrFunc("_returnsFunc", returnsFunc))
	assertNoFailure(t, RegisterArrFunc("_returnsChan", returnsChan))

	_, failure := EvaluateString(`{{ [1]._returnsFunc() }}`, nil)
	assertUnsupportedResult(t, failure, unsupportedFunc)

	_, failure = EvaluateString(`{{ [1]._returnsChan() }}`, nil)
	assertUnsupportedResult(t, failure, unsupportedChan)
}

func TestCustomFunc_UnsupportedIntegerResult(t *testing.T) {
	returnsFunc := func(int, ...any) any { return unsupportedFunc }
	returnsChan := func(int, ...any) any { return unsupportedChan }

	assertNoFailure(t, RegisterIntFunc("_returnsFunc", returnsFunc))
	assertNoFailure(t, RegisterIntFunc("_returnsChan", returnsChan))

	_, failure := EvaluateString(`{{ 1._returnsFunc() }}`, nil)
	assertUnsupportedResult(t, failure, unsupportedFunc)

	_, failure = EvaluateString(`{{ 1._returnsChan() }}`, nil)
	assertUnsupportedResult(t, failure, unsupportedChan)
}

func TestCustomFunc_UnsupportedFloatResult(t *testing.T) {
	returnsFunc := func(float64, ...any) any { return unsupportedFunc }
	returnsChan := func(float64, ...any) any { return unsupportedChan }

	assertNoFailure(t, RegisterFloatFunc("_returnsFunc", returnsFunc))
	assertNoFailure(t, RegisterFloatFunc("_returnsChan", returnsChan))

	_, failure := EvaluateString(`{{ 1.5._returnsFunc() }}`, nil)
	assertUnsupportedResult(t, failure, unsupportedFunc)

	_, failure = EvaluateString(`{{ 1.5._returnsChan() }}`, nil)
	assertUnsupportedResult(t, failure, unsupportedChan)
}

func TestCustomFunc_UnsupportedBooleanResult(t *testing.T) {
	returnsFunc := func(bool, ...any) any { return unsupportedFunc }
	returnsChan := func(bool, ...any) any { return unsupportedChan }

	assertNoFailure(t, RegisterBoolFunc("_returnsFunc", returnsFunc))
	assertNoFailure(t, RegisterBoolFunc("_returnsChan", returnsChan))

	_, failure := EvaluateString(`{{ true._returnsFunc() }}`, nil)
	assertUnsupportedResult(t, failure, unsupportedFunc)

	_, failure = EvaluateString(`{{ true._returnsChan() }}`, nil)
	assertUnsupportedResult(t, failure, unsupportedChan)
}

func TestCustomFunc_UnsupportedObjectResult(t *testing.T) {
	returnsFunc := func(map[string]any, ...any) any { return unsupportedFunc }
	returnsChan := func(map[string]any, ...any) any { return unsupportedChan }

	assertNoFailure(t, RegisterObjFunc("_returnsFunc", returnsFunc))
	assertNoFailure(t, RegisterObjFunc("_returnsChan", returnsChan))

	_, failure := EvaluateString(`{{ {a: 1}._returnsFunc() }}`, nil)
	assertUnsupportedResult(t, failure, unsupportedFunc)

	_, failure = EvaluateString(`{{ {a: 1}._returnsChan() }}`, nil)
	assertUnsupportedResult(t, failure, unsupportedChan)
}

func TestCustomFunc_UnsupportedNestedResult(t *testing.T) {
	nestedFunc := func(string, ...any) any {
		return map[string]any{"nested": unsupportedFunc}
	}

	nestedChan := func(string, ...any) any {
		return map[string]any{"nested": unsupportedChan}
	}

	assertNoFailure(t, RegisterStrFunc("_returnsNestedFunc", nestedFunc))
	assertNoFailure(t, RegisterStrFunc("_returnsNestedChan", nestedChan))

	_, failure := EvaluateString(`{{ "x"._returnsNestedFunc() }}`, nil)
	assertUnsupportedResult(t, failure, unsupportedFunc)

	_, failure = EvaluateString(`{{ "x"._returnsNestedChan() }}`, nil)
	assertUnsupportedResult(t, failure, unsupportedChan)
}

func TestCustomFunc(t *testing.T) {
	t.Run("register for integer receiver", func(t *testing.T) {
		err := RegisterIntFunc("_double", func(num int, args ...any) any {
			return num * 2
		})

		if err != nil {
			t.Fatalf("Error registering function: %s", err)
		}

		assertEvalOutput(t, "{{ 3._double() }}", nil, "6")
	})

	t.Run("register for float receiver", func(t *testing.T) {
		err := RegisterFloatFunc("_double", func(num float64, args ...any) any {
			return num * 2
		})

		if err != nil {
			t.Fatalf("Error registering function: %s", err)
		}

		assertEvalOutput(t, "{{ 3.5._double() }}", nil, "7.0")
	})

	t.Run("register for array receiver", func(t *testing.T) {
		err := RegisterArrFunc("_addNumber", func(arr []any, args ...any) any {
			firstArg := args[0].(int64)
			arr = append(arr, firstArg)
			return arr
		})

		if err != nil {
			t.Fatalf("Error registering function: %s", err)
		}

		assertEvalOutput(t, "{{ [1, 2]._addNumber(3) }}", nil, "1, 2, 3")
	})

	t.Run("register for object receiver", func(t *testing.T) {
		err := RegisterObjFunc("_addProp", func(obj map[string]any, args ...any) any {
			key := args[0].(string)
			value := args[1]
			obj[key] = value
			return obj
		})

		if err != nil {
			t.Fatalf("error registering function: %s", err)
		}

		inp := "{{ obj = {name: 'Anna'}; obj = obj._addProp('age', 25); obj.age }}"
		assertEvalOutput(t, inp, nil, "25")
	})

	t.Run("register for boolean receiver", func(t *testing.T) {
		err := RegisterBoolFunc("_negate", func(b bool, args ...any) any {
			return !b
		})
		if err != nil {
			t.Fatalf("Error registering function: %s", err)
		}

		assertEvalOutput(t, "{{ true._negate() }}", nil, "0")
	})

	t.Run("register for string receiver", func(t *testing.T) {
		err := RegisterStrFunc("_concat", func(s string, args ...any) any {
			arg1Value := args[0].(string)
			arg2Value := args[1].(string)

			return s + arg1Value + arg2Value
		})

		if err != nil {
			t.Fatalf("Error registering function: %s", err)
		}

		assertEvalOutput(t, "{{ 'anna'._concat(' ', 'cho') }}", nil, "anna cho")
	})

	t.Run("registering already registered function", func(t *testing.T) {
		err := RegisterStrFunc("_len", func(s string, args ...any) any {
			return "some output"
		})

		if err != nil {
			t.Fatalf("Error registering function: %s", err)
		}

		// Registering the same function again should return an error
		err = RegisterStrFunc("_len", func(s string, args ...any) any {
			return "some output"
		})

		if err == nil {
			t.Fatalf("Expect error but got none")
		}

		expect := fail.New(nil, "", fail.OriginTpl, fail.ErrFuncAlreadyDefined, "_len", "strings")

		if err := compareFailures(err, expect); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("redefining built-in function not working", func(t *testing.T) {
		err := RegisterStrFunc("trim", func(s string, args ...any) any {
			return "some output"
		})

		if err != nil {
			t.Fatalf("Error registering function: %s", err)
		}

		// the output should be the same as the built-in function
		// even though we redefined it.
		assertEvalOutput(t, "{{ ' anna '.trim() }}", nil, "anna")
	})
}

func TestCustomStrFunc_RawArgArrivesUnescaped(t *testing.T) {
	err := RegisterStrFunc("_returnsInput", func(s string, _ ...any) any {
		return s
	})

	assertNoFailure(t, err)

	// The function gets "<b>", so its result is escaped once
	assertEvalOutput(t, `{{ raw(x)._returnsInput() }}`, map[string]any{"x": "<b>"}, "&lt;b&gt;")

	// Test other cases
	assertEvalOutput(t, `{!! raw(x)._returnsInput() !!}`, map[string]any{"x": "<b>"}, "<b>")
	assertEvalOutput(t, `{!! x._returnsInput() !!}`, map[string]any{"x": "<b>"}, "<b>")
	assertEvalOutput(t, `{{ x._returnsInput() }}`, map[string]any{"x": "<b>"}, "&lt;b&gt;")
}
