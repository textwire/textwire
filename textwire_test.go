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
	var age *int

	cases := []struct {
		name   string
		inp    string
		expect string
		data   map[string]any
	}{
		{
			name:   "Accessing propery name on empty obj variable",
			inp:    `@if(item && item.items && item.items > 1)YES@elseNO@end`,
			expect: "NO",
			data:   map[string]any{"item": nil},
		},
		{
			name:   "Accessing propery name on empty obj variable",
			inp:    `@if(item && item.items && item.items > 1)YES@elseNO@end`,
			expect: "NO",
			data:   map[string]any{"item": nil},
		},
		{
			name:   "Accessing propery name on empty obj variable",
			inp:    `{{ long = user.name.len() > 0 }}{{ long }}`,
			expect: "1",
			data: map[string]any{
				"user": struct{ Name string }{Name: "Harry"},
			},
		},
		{
			name:   "Accessing propery name on empty obj variable",
			inp:    `{{ nameLen = user.name.len() }}{{ nameLen }}`,
			expect: "5",
			data: map[string]any{
				"user": struct{ Name string }{Name: "Harry"},
			},
		},
		{
			name:   "Accessing propery name on empty obj variable",
			inp:    `<p>{{ age }}</p>`,
			expect: `<p></p>`,
			data:   map[string]any{"age": age},
		},
		{
			name:   "Accessing propery name on empty obj variable",
			inp:    `<p>{{ obj = {}; obj.name }}</p>`,
			expect: "<p></p>",
			data:   nil,
		},
		{
			name:   "Accessing property test on empty object {}",
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
		id     int
		inp    string
		expect string
		data   map[string]any
	}{
		{1, `{{ defined('') }}`, "1", nil},
		{2, `{{ defined("") }}`, "1", nil},
		{3, `{{ defined(-0) }}`, "1", nil},
		{3, `{{ defined(0) }}`, "1", nil},
		{4, `{{ defined(1) }}`, "1", nil},
		{5, `{{ defined(0.0) }}`, "1", nil},
		{5, `{{ defined(-0.0) }}`, "1", nil},
		{6, `{{ defined(1.0) }}`, "1", nil},
		{7, `{{ defined({}) }}`, "1", nil},
		{8, `{{ defined([]) }}`, "1", nil},
		{9, `{{ defined(true) }}`, "1", nil},
		{10, `{{ defined(false) }}`, "1", nil},
		{11, `{{ defined(nil) }}`, "1", nil},
		{12, `{{ defined(undefinedVar) }}`, "0", nil},
		{13, `@if(!defined(definedVar))YES@end`, "YES", nil},
		{14, `{{ defined(definedVar) }}`, "1", map[string]any{"definedVar": "nice"}},
		{15, `@if(defined(definedVar))YES@end`, "YES", map[string]any{"definedVar": "nice"}},
		{
			16,
			`{{ defined(definedVar).then("Yes", "No") }}`,
			"Yes",
			map[string]any{"definedVar": "nice"},
		},
		{17, `{{ defined(nilVar) }}`, "1", map[string]any{"nilVar": nil}},
		{18, `@if(defined(nilVar))YES@end`, "YES", map[string]any{"nilVar": nil}},
		{19, `{{ defined(emptyStr) }}`, "1", map[string]any{"emptyStr": ""}},
		{20, `@if(defined(emptyStr))YES@end`, "YES", map[string]any{"emptyStr": ""}},
		{21, `{{ defined(falseVar) }}`, "1", map[string]any{"falseVar": false}},
		{22, `@if(defined(falseVar))YES@end`, "YES", map[string]any{"falseVar": false}},
		{23, `{{ defined(zeroInt) }}`, "1", map[string]any{"zeroInt": 0}},
		{24, `@if(defined(zeroInt))YES@end`, "YES", map[string]any{"zeroInt": 0}},
		{25, `{{ defined(zeroFloat) }}`, "1", map[string]any{"zeroFloat": 0.0}},
		{
			26,
			`{{ defined(obj.prop) }}`,
			"1",
			map[string]any{"obj": map[string]any{"prop": "value"}},
		},
		{27, `{{ defined(obj.prop) }}`, "1", map[string]any{"obj": map[string]any{}}},
		{
			28,
			`{{ defined(obj.prop.test.nice.cool) }}`,
			"0",
			map[string]any{"obj": map[string]any{}},
		},
		{
			29,
			`{{ defined(obj.nested.prop) }}`,
			"1",
			map[string]any{"obj": map[string]any{"nested": map[string]any{"prop": "value"}}},
		},
		{30, `{{ defined(arr[0]) }}`, "1", map[string]any{"arr": []any{"first", "second"}}},
		{
			31,
			`@if(defined(obj.prop))YES@end`,
			"YES",
			map[string]any{"obj": map[string]any{"prop": "value"}},
		},
		{
			32,
			`@if(defined(definedVar) && defined(nilVar))YES@end`,
			"YES",
			map[string]any{"definedVar": "nice", "nilVar": nil},
		},
		{
			33,
			`@if(defined(definedVar) || defined(undefinedVar))YES@end`,
			"YES",
			map[string]any{"definedVar": "nice"},
		},
		{
			34,
			`@if(defined(obj.prop, obj.nested.prop))YES@end`,
			"YES",
			map[string]any{
				"obj": map[string]any{"prop": "value", "nested": map[string]any{"prop": "value"}},
			},
		},
	}

	for _, tc := range cases {
		res, err := EvaluateString(tc.inp, tc.data)
		if err != nil {
			t.Fatalf("Case %d. We don't expect error but got %s", tc.id, err)
		}

		if tc.expect != res {
			t.Errorf("Case %d. Wrong result. Expect: %q got: %q", tc.id, tc.expect, res)
		}
	}
}

func TestHasValueCallExpression(t *testing.T) {
	cases := []struct {
		id     int
		inp    string
		expect string
		data   map[string]any
	}{
		{1, `{{ hasValue('') }}`, "0", nil},
		{2, `{{ hasValue("") }}`, "0", nil},
		{3, `{{ hasValue(0) }}`, "0", nil},
		{3, `{{ hasValue(-0) }}`, "0", nil},
		{4, `{{ hasValue(1) }}`, "1", nil},
		{5, `{{ hasValue(0.0) }}`, "0", nil},
		{5, `{{ hasValue(-0.0) }}`, "0", nil},
		{6, `{{ hasValue(1.0) }}`, "1", nil},
		{7, `{{ hasValue({}) }}`, "0", nil},
		{8, `{{ hasValue([]) }}`, "0", nil},
		{9, `{{ hasValue(true) }}`, "1", nil},
		{10, `{{ hasValue(false) }}`, "0", nil},
		{11, `{{ hasValue(nil) }}`, "0", nil},
		{12, `{{ hasValue(undefinedVar) }}`, "0", nil},
		{13, `@if(!hasValue(definedVar))YES@end`, "YES", nil},
		{14, `{{ hasValue(emptyStr) }}`, "0", map[string]any{"emptyStr": ""}},
		{15, `{{ hasValue(zeroInt) }}`, "0", map[string]any{"zeroInt": 0}},
		{16, `{{ hasValue(zeroFloat) }}`, "0", map[string]any{"zeroFloat": 0.0}},
		{17, `{{ hasValue(falseVar) }}`, "0", map[string]any{"falseVar": false}},
		{18, `{{ hasValue(nilVar) }}`, "0", map[string]any{"nilVar": nil}},
		{19, `{{ hasValue(emptyObj) }}`, "0", map[string]any{"emptyObj": map[string]any{}}},
		{20, `{{ hasValue(emptyArr) }}`, "0", map[string]any{"emptyArr": []any{}}},
		{21, `{{ hasValue(definedVar) }}`, "1", map[string]any{"definedVar": "nice"}},
		{22, `{{ hasValue(nonEmptyStr) }}`, "1", map[string]any{"nonEmptyStr": "hello"}},
		{23, `{{ hasValue(positiveInt) }}`, "1", map[string]any{"positiveInt": 42}},
		{24, `{{ hasValue(positiveFloat) }}`, "1", map[string]any{"positiveFloat": 3.14}},
		{25, `{{ hasValue(trueVar) }}`, "1", map[string]any{"trueVar": true}},
		{
			26,
			`{{ hasValue(nonEmptyObj) }}`,
			"1",
			map[string]any{"nonEmptyObj": map[string]any{"key": "val"}},
		},
		{27, `{{ hasValue(nonEmptyArr) }}`, "1", map[string]any{"nonEmptyArr": []any{1, 2, 3}}},
		{28, `@if(hasValue(definedVar))YES@end`, "YES", map[string]any{"definedVar": "nice"}},
		{29, `@if(hasValue(emptyStr))YES@elseNO@end`, "NO", map[string]any{"emptyStr": ""}},
		{30, `@if(hasValue(nilVar))YES@elseNO@end`, "NO", map[string]any{"nilVar": nil}},
		{31, `@if(hasValue(zeroInt))YES@elseNO@end`, "NO", map[string]any{"zeroInt": 0}},
		{32, `@if(hasValue(falseVar))YES@elseNO@end`, "NO", map[string]any{"falseVar": false}},
		{
			33,
			`{{ hasValue(definedVar).then("Yes", "No") }}`,
			"Yes",
			map[string]any{"definedVar": "nice"},
		},
		{34, `{{ hasValue(emptyStr).then("Yes", "No") }}`, "No", map[string]any{"emptyStr": ""}},
		{
			35,
			`{{ hasValue(obj.prop) }}`,
			"1",
			map[string]any{"obj": map[string]any{"prop": "value"}},
		},
		{36, `{{ hasValue(obj.prop) }}`, "0", map[string]any{"obj": map[string]any{"prop": ""}}},
		{37, `{{ hasValue(obj.prop) }}`, "0", map[string]any{"obj": map[string]any{"prop": 0}}},
		{38, `{{ hasValue(obj.prop) }}`, "0", map[string]any{"obj": map[string]any{}}},
		{39, `{{ hasValue(obj.missing) }}`, "0", map[string]any{"obj": map[string]any{}}},
		{
			40,
			`{{ hasValue(obj.nested.prop) }}`,
			"1",
			map[string]any{"obj": map[string]any{"nested": map[string]any{"prop": "value"}}},
		},
		{
			41,
			`{{ hasValue(obj.nested.prop) }}`,
			"0",
			map[string]any{"obj": map[string]any{"nested": map[string]any{"prop": nil}}},
		},
		{42, `{{ hasValue(arr[0]) }}`, "1", map[string]any{"arr": []any{"first", "second"}}},
		{43, `{{ hasValue(arr[0]) }}`, "0", map[string]any{"arr": []any{""}}},
		{44, `{{ hasValue(arr[0]) }}`, "0", map[string]any{"arr": []any{}}},
		{
			45,
			`@if(defined(definedVar) && hasValue(definedVar))YES@end`,
			"YES",
			map[string]any{"definedVar": "nice"},
		},
		{
			46,
			`@if(defined(definedVar) && hasValue(definedVar))YES@elseNO@end`,
			"NO",
			map[string]any{"definedVar": ""},
		},
		{
			47,
			`@if(hasValue(obj.prop) && hasValue(obj.prop2))YES@end`,
			"YES",
			map[string]any{"obj": map[string]any{"prop": "a", "prop2": "b"}},
		},
		{
			48,
			`@if(hasValue(obj.prop) || hasValue(obj.prop2))YES@end`,
			"YES",
			map[string]any{"obj": map[string]any{"prop": "a"}},
		},
		{49, `{{ hasValue(a, b) }}`, "1", map[string]any{"a": "hello", "b": "world"}},
		{50, `{{ hasValue(a, b) }}`, "0", map[string]any{"a": "hello", "b": ""}},
		{51, `{{ hasValue(a, b) }}`, "0", map[string]any{"a": "", "b": "hello"}},
		{52, `{{ hasValue(a, b) }}`, "0", map[string]any{"a": "", "b": ""}},
		{53, `{{ hasValue(a, b) }}`, "0", map[string]any{"a": 0, "b": 1}},
		{54, `{{ hasValue(a, b) }}`, "0", map[string]any{"a": 1, "b": 0}},
		{55, `{{ hasValue(a, b) }}`, "0", map[string]any{"a": nil, "b": "hello"}},
		{56, `{{ hasValue(a, b, c) }}`, "1", map[string]any{"a": "a", "b": "b", "c": "c"}},
		{57, `{{ hasValue(a, b, c) }}`, "0", map[string]any{"a": "a", "b": "", "c": "c"}},
		{58, `@if(hasValue(a, b))YES@end`, "YES", map[string]any{"a": "x", "b": "y"}},
		{59, `@if(hasValue(a, b))YES@elseNO@end`, "NO", map[string]any{"a": "x", "b": ""}},
		{60, `@if(hasValue(obj.prop, obj.prop2))YES@end`, "YES", map[string]any{"obj": map[string]any{"prop": "a", "prop2": "b"}}},
		{61, `@if(hasValue(obj.prop, obj.prop2))YES@elseNO@end`, "NO", map[string]any{"obj": map[string]any{"prop": "a"}}},
	}

	for _, tc := range cases {
		res, err := EvaluateString(tc.inp, tc.data)
		if err != nil {
			t.Fatalf("Case %d. We don't expect error but got %s", tc.id, err)
		}

		if tc.expect != res {
			t.Errorf("Case %d. Wrong result. Expect: %q got: %q", tc.id, tc.expect, res)
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
			inp: `{{ defined(name.undefinedFunc()) }}`,
			err: fail.New(
				1,
				"",
				"evaluator",
				fail.ErrFuncNotDefined,
				object.STR_OBJ,
				"undefinedFunc",
			),
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
