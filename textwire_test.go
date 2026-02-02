package textwire

import (
	"io"
	"log"
	"os"
	"testing"

	"github.com/textwire/textwire/v3/fail"
	"github.com/textwire/textwire/v3/object"
)

func readFile(fileName string) (string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return "", err
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("readFile: failed to close %s: %v", fileName, err)
		}
	}()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func TestEvaluateString(t *testing.T) {
	cases := []struct {
		name   string
		inp    string
		expect string
		data   map[string]any
	}{
		{
			name:   "Simple math operation with integers",
			inp:    "{{ 1 + 2 }}",
			expect: "3",
			data:   nil,
		},
		{
			name:   "Simple math operation with identifiers",
			inp:    "{{ n1 * n2 }}",
			expect: "2",
			data:   map[string]any{"n1": 1, "n2": 2},
		},
		{
			name:   "First letter of the object property is case insensitive",
			inp:    "{{ user.iD.str() + ' ' + user.ID.str() }}",
			expect: "1 1",
			data: map[string]any{
				"user": struct{ ID uint }{1}},
		},
		{
			name:   "Accessing user.name.firstName property",
			inp:    "{{ user.name.firstName }}",
			expect: "Ann",
			data: map[string]any{
				"user": struct {
					Name struct{ FirstName string }
				}{
					Name: struct{ FirstName string }{"Ann"},
				},
			},
		},
		{
			name:   "Empty global object is defined",
			inp:    "<span>{{ global }}</span>",
			expect: "<span>{}</span>",
			data:   nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := EvaluateString(tc.inp, tc.data)
			if err != nil {
				t.Errorf("error evaluating template: %s", err)
			}

			if actual != tc.expect {
				t.Errorf("wrong result. expect:\n\"%s\"\ngot:\n\"%s\"",
					tc.expect, actual)
			}
		})
	}
}

func TestDefinedCallExpression(t *testing.T) {
	cases := []struct {
		inp    string
		expect string
		data   map[string]any
	}{
		{`{{ defined('') }}`, "1", nil},
		{`{{ defined("") }}`, "1", nil},
		{`{{ defined(0) }}`, "1", nil},
		{`{{ defined(1) }}`, "1", nil},
		{`{{ defined(0.0) }}`, "1", nil},
		{`{{ defined(1.0) }}`, "1", nil},
		{`{{ defined({}) }}`, "1", nil},
		{`{{ defined([]) }}`, "1", nil},
		{`{{ defined(true) }}`, "1", nil},
		{`{{ defined(false) }}`, "1", nil},
		{`{{ defined(nil) }}`, "1", nil},
		{`{{ defined(undefinedVar) }}`, "0", nil},
		{`@if(!defined(definedVar))YES@end`, "YES", nil},
		{
			inp:    `{{ defined(definedVar) }}`,
			expect: "1",
			data:   map[string]any{"definedVar": "nice"},
		},
		{
			inp:    `@if(defined(definedVar))YES@end`,
			expect: "YES",
			data:   map[string]any{"definedVar": "nice"},
		},
		{
			inp:    `{{ defined(definedVar).then("Yes", "No") }}`,
			expect: "Yes",
			data:   map[string]any{"definedVar": "nice"}},

		// Variables with falsy but defined values like nil, false, ""
		{`{{ defined(nilVar) }}`, "1", map[string]any{"nilVar": nil}},
		{`@if(defined(nilVar))YES@end`, "YES", map[string]any{"nilVar": nil}},
		{`{{ defined(emptyStr) }}`, "1", map[string]any{"emptyStr": ""}},
		{`@if(defined(emptyStr))YES@end`, "YES", map[string]any{"emptyStr": ""}},
		{`{{ defined(falseVar) }}`, "1", map[string]any{"falseVar": false}},
		{`@if(defined(falseVar))YES@end`, "YES", map[string]any{"falseVar": false}},
		{`{{ defined(zeroInt) }}`, "1", map[string]any{"zeroInt": 0}},
		{`@if(defined(zeroInt))YES@end`, "YES", map[string]any{"zeroInt": 0}},
		{`{{ defined(zeroFloat) }}`, "1", map[string]any{"zeroFloat": 0.0}},

		// Complex data structures with nested objects
		{
			inp:    `{{ defined(obj.prop) }}`,
			expect: "1",
			data:   map[string]any{"obj": map[string]any{"prop": "value"}},
		},
		{
			inp:    `{{ defined(obj.nested.prop) }}`,
			expect: "1",
			data: map[string]any{
				"obj": map[string]any{
					"nested": map[string]any{"prop": "value"},
				},
			},
		},
		{
			inp:    `{{ defined(arr[0]) }}`,
			expect: "1",
			data:   map[string]any{"arr": []any{"first", "second"}}},

		// More conditional logic tests
		{
			inp:    `@if(defined(obj.prop))YES@end`,
			expect: "YES",
			data:   map[string]any{"obj": map[string]any{"prop": "value"}},
		},
		{
			inp:    `@if(defined(definedVar) && defined(nilVar))YES@end`,
			expect: "YES",
			data:   map[string]any{"definedVar": "nice", "nilVar": nil},
		},
		{
			inp:    `@if(defined(definedVar) || defined(undefinedVar))YES@end`,
			expect: "YES",
			data:   map[string]any{"definedVar": "nice"},
		},
		{
			inp:    `@if(defined(obj.prop, obj.nested.prop))YES@end`,
			expect: "YES",
			data: map[string]any{
				"obj": map[string]any{
					"prop":   "value",
					"nested": map[string]any{"prop": "value"},
				},
			},
		},
	}

	for _, tc := range cases {
		res, err := EvaluateString(tc.inp, tc.data)
		if err != nil {
			t.Fatalf("we don't expect error but got %s", err)
		}

		if tc.expect != res {
			t.Errorf("wrong result. expect: %q got: %q", tc.expect, res)
		}
	}
}

func TestErrorHandling(t *testing.T) {
	cases := []struct {
		inp  string
		err  *fail.Error
		data map[string]any
	}{
		{
			inp:  `@use("someTemplate")`,
			err:  fail.New(1, "", "evaluator", fail.ErrSomeDirsOnlyInTemplates),
			data: nil,
		},
		{
			inp:  `@insert("title", "hi")`,
			err:  fail.New(1, "", "evaluator", fail.ErrSomeDirsOnlyInTemplates),
			data: nil,
		},
		{
			inp:  `@reserve("content")`,
			err:  fail.New(1, "", "evaluator", fail.ErrSomeDirsOnlyInTemplates),
			data: nil,
		},
		{
			inp:  `@component("~header")`,
			err:  fail.New(1, "", "evaluator", fail.ErrSomeDirsOnlyInTemplates),
			data: nil,
		},
		{
			inp: `{{ 1 + "a" }}`,
			err: fail.New(
				1,
				"",
				"evaluator",
				fail.ErrTypeMismatch,
				object.INT_OBJ,
				"+",
				object.STR_OBJ,
			),
			data: nil,
		},
		{
			inp:  `{{ loop = "test" }}`,
			err:  fail.New(1, "", "evaluator", fail.ErrReservedIdentifiers),
			data: nil,
		},
		{
			inp:  `{{ global = "test" }}`,
			err:  fail.New(1, "", "evaluator", fail.ErrReservedIdentifiers),
			data: nil,
		},
		{
			inp: `{{ loop }}`,
			err: fail.New(
				0,
				"",
				"evaluator",
				fail.ErrReservedIdentifiers,
			), data: map[string]any{"loop": "test"},
		},
		{
			inp: `{{ n = 1; n = "test" }}`,
			err: fail.New(
				1,
				"",
				"evaluator",
				fail.ErrIdentifierTypeMismatch,
				"n",
				object.INT_OBJ,
				object.STR_OBJ,
			),
			data: nil,
		},
		{
			inp:  `{{ obj = {}; obj.name }}`,
			err:  fail.New(1, "", "evaluator", fail.ErrPropertyNotFound, "name", object.OBJ_OBJ),
			data: nil,
		},
		{
			inp:  `{{ {}.test }}`,
			err:  fail.New(1, "", "evaluator", fail.ErrPropertyNotFound, "test", object.OBJ_OBJ),
			data: nil,
		},
		{
			inp: `{{ 5.someFunction() }}`,
			err: fail.New(
				1,
				"",
				"evaluator",
				fail.ErrNoFuncForThisType,
				"someFunction",
				object.INT_OBJ,
			),
			data: nil,
		},
		{
			inp:  `{{ 3 / 0 }}`,
			err:  fail.New(1, "", "evaluator", fail.ErrDivisionByZero),
			data: nil,
		},
		{
			inp:  `{{ undefinedVar }}`,
			err:  fail.New(1, "", "evaluator", fail.ErrIdentifierIsUndefined, "undefinedVar"),
			data: nil,
		},
		{
			inp:  `{{ obj = {name: "Amy"}; obj.name.id }}`,
			err:  fail.New(1, "", "evaluator", fail.ErrProperyOnNonObject, "id", object.STR_OBJ),
			data: nil,
		},
	}

	for _, tc := range cases {
		_, err := EvaluateString(tc.inp, tc.data)
		if err == nil {
			t.Fatalf("expect error but got none")
		}

		if err.Error() != tc.err.String() {
			t.Errorf("wrong error message. expected:\n%q\ngot:\n%q",
				tc.err, err)
		}
	}
}

func TestEvaluateFile(t *testing.T) {
	filename := "14.two-vars-no-layout"
	absPath, fileErr := getFullPath("textwire/testdata/good/before/" + filename + ".tw")

	if fileErr != nil {
		t.Errorf("error getting full path: %s", fileErr)
		return
	}

	out, err := EvaluateFile(absPath, map[string]any{
		"title":   "Textwire is Awesome",
		"visible": true,
	})

	if err != nil {
		t.Errorf("error evaluating file:\n%s", err)
	}

	expect, err := readFile("textwire/testdata/good/expected/" + filename + ".html")
	if err != nil {
		t.Errorf("error reading expected file: %s", err)
		return
	}

	if out != expect {
		t.Errorf("wrong output. expect:\n%s\ngot:\n%s", expect, out)
	}
}

func TestCustomFunctions(t *testing.T) {
	t.Run("register for integer receiver", func(t *testing.T) {
		err := RegisterIntFunc("_double", func(num int, args ...any) any {
			return num * 2
		})

		if err != nil {
			t.Fatalf("error registering function: %s", err)
		}

		actual, err := EvaluateString("{{ 3._double() }}", nil)
		if err != nil {
			t.Fatalf("error evaluating template: %s", err)
		}

		if actual != "6" {
			t.Errorf("wrong result. expect: '6' got: '%s'", actual)
		}
	})

	t.Run("register for float receiver", func(t *testing.T) {
		err := RegisterFloatFunc("_double", func(num float64, args ...any) any {
			return num * 2
		})

		if err != nil {
			t.Fatalf("error registering function: %s", err)
		}

		actual, err := EvaluateString("{{ 3.5._double() }}", nil)
		if err != nil {
			t.Fatalf("error evaluating template: %s", err)
		}

		if actual != "7.0" {
			t.Fatalf("wrong result. expect: '7.0' got: '%s'", actual)
		}
	})

	t.Run("register for array receiver", func(t *testing.T) {
		err := RegisterArrFunc("_addNumber", func(arr []any, args ...any) any {
			firstArg := args[0].(int64)
			arr = append(arr, firstArg)
			return arr
		})

		if err != nil {
			t.Fatalf("error registering function: %s", err)
		}

		actual, err := EvaluateString("{{ [1, 2]._addNumber(3) }}", nil)
		if err != nil {
			t.Fatalf("error evaluating template: %s", err)
		}

		if actual != "1, 2, 3" {
			t.Fatalf("wrong result. expect: '1, 2, 3' got: '%s'", actual)
		}
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
		actual, err := EvaluateString(inp, nil)
		if err != nil {
			t.Fatalf("error evaluating template: %s", err)
		}

		if actual != "25" {
			t.Fatalf("wrong result. expect: '25' got: '%s'", actual)
		}
	})

	t.Run("register for boolean receiver", func(t *testing.T) {
		err := RegisterBoolFunc("_negate", func(b bool, args ...any) any {
			return !b
		})
		if err != nil {
			t.Fatalf("error registering function: %s", err)
		}

		actual, err := EvaluateString("{{ true._negate() }}", nil)
		if err != nil {
			t.Fatalf("error evaluating template: %s", err)
		}

		if actual != "0" {
			t.Fatalf("wrong result. expect: '0' got '%s'", actual)
		}
	})

	t.Run("register for string receiver", func(t *testing.T) {
		err := RegisterStrFunc("_concat", func(s string, args ...any) any {
			arg1Value := args[0].(string)
			arg2Value := args[1].(string)

			return s + arg1Value + arg2Value
		})

		if err != nil {
			t.Fatalf("error registering function: %s", err)
		}

		actual, err := EvaluateString("{{ 'anna'._concat(' ', 'cho') }}", nil)
		if err != nil {
			t.Fatalf("error evaluating template: %s", err)
		}

		if actual != "anna cho" {
			t.Fatalf("wrong result. expect: 'anna cho' got: '%s'", actual)
		}
	})

	t.Run("registering already registered function", func(t *testing.T) {
		err := RegisterStrFunc("_len", func(s string, args ...any) any {
			return "some output"
		})

		if err != nil {
			t.Fatalf("error registering function: %s", err)
		}

		// Registering the same function again should return an error
		err = RegisterStrFunc("_len", func(s string, args ...any) any {
			return "some output"
		})

		if err == nil {
			t.Fatalf("expect error but got none")
		}

		expect := fail.New(0, "", "API", fail.ErrFuncAlreadyDefined,
			"_len", "strings")

		if err.Error() != expect.Error().Error() {
			t.Fatalf("wrong error message. expect: %q got: %q", expect, err)
		}
	})

	t.Run("redefining built-in function not working", func(t *testing.T) {
		err := RegisterStrFunc("trim", func(s string, args ...any) any {
			return "some output"
		})

		if err != nil {
			t.Fatalf("error registering function: %s", err)
		}

		actual, err := EvaluateString("{{ ' anna '.trim() }}", nil)
		if err != nil {
			t.Fatalf("error registering function: %s", err)
		}

		// the output should be the same as the built-in function
		// even though we redefined it.
		if actual != "anna" {
			t.Fatalf("wrong output. expect: 'anna' got: '%s'", actual)
		}
	})
}
