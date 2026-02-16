package textwire

import (
	"io"
	"log"
	"os"
	"testing"

	"github.com/textwire/textwire/v3/pkg/fail"
	"github.com/textwire/textwire/v3/pkg/file"
	"github.com/textwire/textwire/v3/pkg/object"
	"github.com/textwire/textwire/v3/pkg/token"
)

func readFile(fileName string) (string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return "", err
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Function readFile failed to close %s: %v", fileName, err)
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
			name:   "Accessing propery 'name' on empty 'obj' variable",
			inp:    `<p>{{ obj = {}; obj.name }}</p>`,
			expect: "<p></p>",
			data:   nil,
		},
		{
			name:   "Accessing property 'test' on empty object '{}'",
			inp:    `<h2>{{ {}.test }}</h2>`,
			expect: "<h2></h2>",
			data:   nil,
		},
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
				t.Errorf("Wrong evaluation result. Expect:\n'%s'\ngot:\n'%s'", tc.expect, actual)
			}
		})
	}
}

func TestDefinedCallExpression(t *testing.T) {
	cases := []struct {
		name   string
		inp    string
		expect string
		data   map[string]any
	}{
		{
			name:   "defined on empty single quoted string literal",
			inp:    `{{ defined('') }}`,
			expect: "1",
			data:   nil,
		},
		{
			name:   "defined on empty double quoted string literal",
			inp:    `{{ defined("") }}`,
			expect: "1",
			data:   nil,
		},
		{name: "defined on zero int literal", inp: `{{ defined(0) }}`, expect: "1", data: nil},
		{name: "defined on one int literal", inp: `{{ defined(1) }}`, expect: "1", data: nil},
		{name: "defined on zero float literal", inp: `{{ defined(0.0) }}`, expect: "1", data: nil},
		{name: "defined on one float literal", inp: `{{ defined(1.0) }}`, expect: "1", data: nil},
		{name: "defined on empty object literal", inp: `{{ defined({}) }}`, expect: "1", data: nil},
		{name: "defined on empty array literal", inp: `{{ defined([]) }}`, expect: "1", data: nil},
		{name: "defined on true bool literal", inp: `{{ defined(true) }}`, expect: "1", data: nil},
		{
			name:   "defined on false bool literal",
			inp:    `{{ defined(false) }}`,
			expect: "1",
			data:   nil,
		},
		{name: "defined on nil literal", inp: `{{ defined(nil) }}`, expect: "1", data: nil},
		{
			name:   "defined on undefined variable returns false",
			inp:    `{{ defined(undefinedVar) }}`,
			expect: "0",
			data:   nil,
		},
		{
			name:   "if not defined on undefined variable",
			inp:    `@if(!defined(definedVar))YES@end`,
			expect: "YES",
			data:   nil,
		},
		{
			name:   "defined on existing string variable",
			inp:    `{{ defined(definedVar) }}`,
			expect: "1",
			data:   map[string]any{"definedVar": "nice"},
		},
		{
			name:   "if defined on existing string variable",
			inp:    `@if(defined(definedVar))YES@end`,
			expect: "YES",
			data:   map[string]any{"definedVar": "nice"},
		},
		{
			name:   "defined result with then method",
			inp:    `{{ defined(definedVar).then("Yes", "No") }}`,
			expect: "Yes",
			data:   map[string]any{"definedVar": "nice"},
		},
		{
			name:   "defined on nil variable",
			inp:    `{{ defined(nilVar) }}`,
			expect: "1",
			data:   map[string]any{"nilVar": nil},
		},
		{
			name:   "if defined on nil variable",
			inp:    `@if(defined(nilVar))YES@end`,
			expect: "YES",
			data:   map[string]any{"nilVar": nil},
		},
		{
			name:   "defined on empty string variable",
			inp:    `{{ defined(emptyStr) }}`,
			expect: "1",
			data:   map[string]any{"emptyStr": ""},
		},
		{
			name:   "if defined on empty string variable",
			inp:    `@if(defined(emptyStr))YES@end`,
			expect: "YES",
			data:   map[string]any{"emptyStr": ""},
		},
		{
			name:   "defined on false variable",
			inp:    `{{ defined(falseVar) }}`,
			expect: "1",
			data:   map[string]any{"falseVar": false},
		},
		{
			name:   "if defined on false variable",
			inp:    `@if(defined(falseVar))YES@end`,
			expect: "YES",
			data:   map[string]any{"falseVar": false},
		},
		{
			name:   "defined on zero int variable",
			inp:    `{{ defined(zeroInt) }}`,
			expect: "1",
			data:   map[string]any{"zeroInt": 0},
		},
		{
			name:   "if defined on zero int variable",
			inp:    `@if(defined(zeroInt))YES@end`,
			expect: "YES",
			data:   map[string]any{"zeroInt": 0},
		},
		{
			name:   "defined on zero float variable",
			inp:    `{{ defined(zeroFloat) }}`,
			expect: "1",
			data:   map[string]any{"zeroFloat": 0.0},
		},
		{
			name:   "defined on existing object property",
			inp:    `{{ defined(obj.prop) }}`,
			expect: "1",
			data:   map[string]any{"obj": map[string]any{"prop": "value"}},
		},
		{
			name:   "defined on missing object property",
			inp:    `{{ defined(obj.prop) }}`,
			expect: "1",
			data:   map[string]any{"obj": map[string]any{}},
		},
		{
			name:   "defined on missing object property",
			inp:    `{{ defined(obj.prop.test.nice.cool) }}`,
			expect: "0",
			data:   map[string]any{"obj": map[string]any{}},
		},
		{
			name:   "defined on nested object property",
			inp:    `{{ defined(obj.nested.prop) }}`,
			expect: "1",
			data: map[string]any{
				"obj": map[string]any{
					"nested": map[string]any{"prop": "value"},
				},
			},
		},
		{
			name:   "defined on existing array element",
			inp:    `{{ defined(arr[0]) }}`,
			expect: "1",
			data:   map[string]any{"arr": []any{"first", "second"}},
		},
		{
			name:   "if defined on object property",
			inp:    `@if(defined(obj.prop))YES@end`,
			expect: "YES",
			data:   map[string]any{"obj": map[string]any{"prop": "value"}},
		},
		{
			name:   "defined and operator with defined and nil",
			inp:    `@if(defined(definedVar) && defined(nilVar))YES@end`,
			expect: "YES",
			data:   map[string]any{"definedVar": "nice", "nilVar": nil},
		},
		{
			name:   "defined or operator with defined and undefined",
			inp:    `@if(defined(definedVar) || defined(undefinedVar))YES@end`,
			expect: "YES",
			data:   map[string]any{"definedVar": "nice"},
		},
		{
			name:   "defined with multiple properties",
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
		t.Run(tc.name, func(t *testing.T) {
			res, err := EvaluateString(tc.inp, tc.data)
			if err != nil {
				t.Fatalf("We don't expect error but got %s", err)
			}

			if tc.expect != res {
				t.Errorf("Wrong result. Expect: %q got: %q", tc.expect, res)
			}
		})
	}
}

func TestErrorHandling(t *testing.T) {
	cases := []struct {
		inp  string
		err  *fail.Error
		data map[string]any
	}{
		{
			inp:  `{{ defined(name.undefinedFunc()) }}`,
			err:  fail.New(1, "", "evaluator", fail.ErrFuncNotDefined, object.STR_OBJ, "undefinedFunc"),
			data: map[string]any{"name": "Anna"},
		},
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
			inp:  `{{ user = {}; user.address.zip }}`,
			err:  fail.New(1, "", "evaluator", fail.ErrPropertyOnNonObject, object.NIL_OBJ, "zip"),
			data: nil,
		},
		{
			inp: `{{ 5.someFunction() }}`,
			err: fail.New(
				1,
				"",
				"evaluator",
				fail.ErrFuncNotDefined,
				object.INT_OBJ,
				"someFunction",
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
			err:  fail.New(1, "", "evaluator", fail.ErrVariableIsUndefined, "undefinedVar"),
			data: nil,
		},
		{
			inp:  `{{ obj = {name: "Amy"}; obj.name.id }}`,
			err:  fail.New(1, "", "evaluator", fail.ErrPropertyOnNonObject, object.STR_OBJ, "id"),
			data: nil,
		},
		{
			inp: `{{ obj."str" }}`,
			err: fail.New(
				1,
				"",
				"evaluator",
				fail.ErrWrongNextToken,
				token.String(token.IDENT),
				token.String(token.STR),
			),
			data: nil,
		},
		{
			inp: `@each(v in {}){{ v }}@end`,
			err: fail.New(
				1,
				"",
				"evaluator",
				fail.ErrEachDirWithNonArrArg,
				object.OBJ_OBJ,
			),
			data: nil,
		},
	}

	for _, tc := range cases {
		_, err := EvaluateString(tc.inp, tc.data)
		if err == nil {
			t.Fatalf("Expected error but got none")
		}

		if err.Error() != tc.err.String() {
			t.Errorf("Wrong error message. Expected:\n%q\ngot:\n%q", tc.err, err)
		}
	}
}

func TestEvaluateFile(t *testing.T) {
	absPath, fileErr := file.ToFullPath("testdata/good/before/two-vars-no-use/index.tw")
	if fileErr != nil {
		t.Errorf("Error getting full path: %s", fileErr)
		return
	}

	actual, err := EvaluateFile(absPath, map[string]any{
		"title":   "Textwire is Awesome",
		"visible": true,
	})

	if err != nil {
		t.Errorf("Error evaluating file: %q", err)
	}

	expect, err := readFile("testdata/good/expected/two-vars-no-use.html")
	if err != nil {
		t.Errorf("Error reading file: %q", err)
		return
	}

	if actual != expect {
		t.Errorf("Wrong output. Expect:\n%q\ngot:\n%q", expect, actual)
	}
}

func TestCustomFunctions(t *testing.T) {
	t.Run("register for integer receiver", func(t *testing.T) {
		err := RegisterIntFunc("_double", func(num int, args ...any) any {
			return num * 2
		})

		if err != nil {
			t.Fatalf("Error registering function: %s", err)
		}

		actual, err := EvaluateString("{{ 3._double() }}", nil)
		if err != nil {
			t.Fatalf("Error evaluating template: %s", err)
		}

		if actual != "6" {
			t.Errorf("Wrong result. Expect 6 but got %s", actual)
		}
	})

	t.Run("register for float receiver", func(t *testing.T) {
		err := RegisterFloatFunc("_double", func(num float64, args ...any) any {
			return num * 2
		})

		if err != nil {
			t.Fatalf("Error registering function: %s", err)
		}

		actual, err := EvaluateString("{{ 3.5._double() }}", nil)
		if err != nil {
			t.Fatalf("Error evaluating template: %s", err)
		}

		if actual != "7.0" {
			t.Fatalf("Wrong result. Expect 7.0 but got %s", actual)
		}
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

		actual, err := EvaluateString("{{ [1, 2]._addNumber(3) }}", nil)
		if err != nil {
			t.Fatalf("Error evaluating template: %s", err)
		}

		if actual != "1, 2, 3" {
			t.Fatalf("Wrong result. Expect '1, 2, 3' got '%s'", actual)
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
			t.Fatalf("Error evaluating template: %s", err)
		}

		if actual != "25" {
			t.Fatalf("Wrong result. Expect 25 but got %s", actual)
		}
	})

	t.Run("register for boolean receiver", func(t *testing.T) {
		err := RegisterBoolFunc("_negate", func(b bool, args ...any) any {
			return !b
		})
		if err != nil {
			t.Fatalf("Error registering function: %s", err)
		}

		actual, err := EvaluateString("{{ true._negate() }}", nil)
		if err != nil {
			t.Fatalf("Error evaluating template: %s", err)
		}

		if actual != "0" {
			t.Fatalf("Wrong result. Expect 0 but got %s", actual)
		}
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

		actual, err := EvaluateString("{{ 'anna'._concat(' ', 'cho') }}", nil)
		if err != nil {
			t.Fatalf("Error evaluating template: %s", err)
		}

		if actual != "anna cho" {
			t.Fatalf("Wrong result. Expect 'anna cho' but got '%s'", actual)
		}
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

		expect := fail.New(0, "", "API", fail.ErrFuncAlreadyDefined, "_len", "strings")
		if err.Error() != expect.Error().Error() {
			t.Fatalf("Wrong error message. Expect:\n%q\ngot:\n%q", expect, err)
		}
	})

	t.Run("redefining built-in function not working", func(t *testing.T) {
		err := RegisterStrFunc("trim", func(s string, args ...any) any {
			return "some output"
		})

		if err != nil {
			t.Fatalf("Error registering function: %s", err)
		}

		actual, err := EvaluateString("{{ ' anna '.trim() }}", nil)
		if err != nil {
			t.Fatalf("Error registering function: %s", err)
		}

		// the output should be the same as the built-in function
		// even though we redefined it.
		if actual != "anna" {
			t.Fatalf("Wrong output. Expect 'anna' but got '%s'", actual)
		}
	})
}
