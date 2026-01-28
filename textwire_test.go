package textwire

import (
	"io"
	"log"
	"os"
	"testing"

	"github.com/textwire/textwire/v2/config"
	"github.com/textwire/textwire/v2/fail"
	"github.com/textwire/textwire/v2/object"
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
		config *config.Config
	}{
		{
			name:   "Simple math operation with integers",
			inp:    "{{ 1 + 2 }}",
			expect: "3",
			data:   nil,
			config: nil,
		},
		{
			name:   "Simple math operation with identifiers",
			inp:    "{{ n1 * n2 }}",
			expect: "2",
			data:   map[string]any{"n1": 1, "n2": 2},
			config: nil,
		},
		{
			name:   "Accessing user.name.firstName property",
			inp:    "{{ user.name.firstName }}",
			expect: "Ann",
			data: map[string]any{
				"user": struct {
					Name struct{ FirstName string }
					Age  int
				}{
					Name: struct{ FirstName string }{"Ann"},
					Age:  20,
				},
			},
			config: nil,
		},
		{
			name:   "Empty global object is defined",
			inp:    "<span>{{ global }}</span>",
			expect: "<span>{}</span>",
			data:   nil,
			config: nil,
		},
		{
			name:   "Getting 'env' value from global object",
			inp:    "<span>{{ global.env }}</span>",
			expect: "<span>development</span>",
			data:   nil,
			config: &config.Config{
				GlobalData: map[string]any{
					"env": "development",
				},
			},
		},
		{
			name:   "Getting multiple values from global object",
			inp:    "<h1>{{ global.first + ' ' + global.last }}</h1>",
			expect: "<h1>Amy Adams</h1>",
			data:   nil,
			config: &config.Config{
				GlobalData: map[string]any{
					"first": "Amy",
					"last":  "Adams",
				},
			},
		},
		{
			name:   "Accessing string array in global object by index",
			inp:    "<span>{{ global.tags[0] }}</span>",
			expect: "<span>go</span>",
			data:   nil,
			config: &config.Config{
				GlobalData: map[string]any{
					"tags": []string{"go", "programming", "templates"},
				},
			},
		},
		{
			name:   "Accessing string array in global object by index",
			inp:    "@if(myVar.binary())NICE@end",
			expect: "NICE",
			data:   map[string]any{"myVar": true},
			config: nil,
		},
		{
			name:   "Accessing string array in global object by index",
			inp:    "@if(myVar.binary())NICE@end",
			expect: "",
			data:   map[string]any{"myVar": false},
			config: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.config != nil {
				Configure(tc.config)
			}

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

func TestIsDefinedCallExpression(t *testing.T) {
	cases := []struct {
		inp    string
		expect string
		data   map[string]any
	}{
		{`{{ definedVar }}`, "nice", map[string]any{"definedVar": "nice"}},
		{`{{ definedVar.isDefined() }}`, "1", map[string]any{"definedVar": "nice"}},
		{`{{ undefinedVar.isDefined() }}`, "0", nil},
		{`@if(definedVar.isDefined())YES@end`, "YES", map[string]any{"definedVar": "nice"}},
		{`@if(!definedVar.isDefined())YES@end`, "YES", nil},

		// Variables with falsy but defined values like nil, false, ""
		{`{{ nilVar.isDefined() }}`, "1", map[string]any{"nilVar": nil}},
		{`@if(nilVar.isDefined())YES@end`, "YES", map[string]any{"nilVar": nil}},
		{`{{ emptyStr.isDefined() }}`, "1", map[string]any{"emptyStr": ""}},
		{`@if(emptyStr.isDefined())YES@end`, "YES", map[string]any{"emptyStr": ""}},
		{`{{ falseVar.isDefined() }}`, "1", map[string]any{"falseVar": false}},
		{`@if(falseVar.isDefined())YES@end`, "YES", map[string]any{"falseVar": false}},
		{`{{ zeroInt.isDefined() }}`, "1", map[string]any{"zeroInt": 0}},
		{`@if(zeroInt.isDefined())YES@end`, "YES", map[string]any{"zeroInt": 0}},
		{`{{ zeroFloat.isDefined() }}`, "1", map[string]any{"zeroFloat": 0.0}},

		// Complex data structures with nested objects
		{
			inp:    `{{ obj.prop.isDefined() }}`,
			expect: "1",
			data:   map[string]any{"obj": map[string]any{"prop": "value"}},
		},
		{
			inp:    `{{ obj.nested.prop.isDefined() }}`,
			expect: "1",
			data: map[string]any{
				"obj": map[string]any{
					"nested": map[string]any{"prop": "value"},
				},
			},
		},
		{
			inp:    `{{ arr[0].isDefined() }}`,
			expect: "1",
			data:   map[string]any{"arr": []any{"first", "second"}}},

		// More conditional logic tests
		{
			inp:    `@if(obj.prop.isDefined())YES@end`,
			expect: "YES",
			data:   map[string]any{"obj": map[string]any{"prop": "value"}},
		},
		{
			inp:    `@if(definedVar.isDefined() && nilVar.isDefined())YES@end`,
			expect: "YES",
			data:   map[string]any{"definedVar": "nice", "nilVar": nil},
		},
		{
			inp:    `@if(definedVar.isDefined() || undefinedVar.isDefined())YES@end`,
			expect: "YES",
			data:   map[string]any{"definedVar": "nice"},
		},
		{
			inp:    `@if(obj.prop.isDefined() && obj.nested.prop.isDefined())YES@end`,
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
		{`{{ 1 + "a" }}`, fail.New(1, "", "evaluator", fail.ErrTypeMismatch, object.INT_OBJ, "+", object.STR_OBJ), nil},
		{`@use("someTemplate")`, fail.New(1, "", "evaluator", fail.ErrUseStmtMissingLayout), nil},
		{`{{ loop = "test" }}`, fail.New(1, "", "evaluator", fail.ErrReservedIdentifiers), nil},
		{`{{ global = "test" }}`, fail.New(1, "", "evaluator", fail.ErrReservedIdentifiers), nil},
		{`{{ loop }}`, fail.New(0, "", "evaluator", fail.ErrReservedIdentifiers), map[string]any{"loop": "test"}},
		{`{{ n = 1; n = "test" }}`, fail.New(1, "", "evaluator", fail.ErrIdentifierTypeMismatch, "n", object.INT_OBJ, object.STR_OBJ), nil},
		{`{{ global.username }}`, fail.New(1, "", "evaluator", fail.ErrPropertyNotFound, "username", object.OBJ_OBJ), nil},
		{`{{ obj = {}; obj.name }}`, fail.New(1, "", "evaluator", fail.ErrPropertyNotFound, "name", object.OBJ_OBJ), nil},
		{`{{ {}.test }}`, fail.New(1, "", "evaluator", fail.ErrPropertyNotFound, "test", object.OBJ_OBJ), nil},
		{`{{ 5.someFunction() }}`, fail.New(1, "", "evaluator", fail.ErrNoFuncForThisType, "someFunction", object.INT_OBJ), nil},
		{`{{ 3 / 0 }}`, fail.New(1, "", "evaluator", fail.ErrDivisionByZero), nil},
		{`{{ 1 ~ 8 }}`, fail.New(1, "", "parser", fail.ErrIllegalToken, "~"), nil},
		{`{{ undefinedVar }}`, fail.New(1, "", "parser", fail.ErrIdentifierIsUndefined, "undefinedVar"), nil},
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
	absPath, fileErr := getFullPath("textwire/testdata/good/before/"+filename+".tw", false)

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
		err := RegisterIntFunc("double", func(num int, args ...any) int {
			return num * 2
		})

		if err != nil {
			t.Fatalf("error registering function: %s", err)
		}

		actual, err := EvaluateString("{{ 3.double() }}", nil)
		if err != nil {
			t.Fatalf("error evaluating template: %s", err)
		}

		if actual != "6" {
			t.Errorf("wrong result. expect: '6' got: '%s'", actual)
		}
	})

	t.Run("register for float receiver", func(t *testing.T) {
		err := RegisterFloatFunc("double", func(num float64, args ...any) float64 {
			return num * 2
		})

		if err != nil {
			t.Fatalf("error registering function: %s", err)
		}

		actual, err := EvaluateString("{{ 3.5.double() }}", nil)
		if err != nil {
			t.Fatalf("error evaluating template: %s", err)
		}

		if actual != "7.0" {
			t.Fatalf("wrong result. expect: '7.0' got: '%s'", actual)
		}
	})

	t.Run("register for array receiver", func(t *testing.T) {
		err := RegisterArrFunc("addNumber", func(arr []any, args ...any) []any {
			firstArg := args[0].(int64)
			arr = append(arr, firstArg)
			return arr
		})

		if err != nil {
			t.Fatalf("error registering function: %s", err)
		}

		actual, err := EvaluateString("{{ [1, 2].addNumber(3) }}", nil)
		if err != nil {
			t.Fatalf("error evaluating template: %s", err)
		}

		if actual != "1, 2, 3" {
			t.Fatalf("wrong result. expect: '1, 2, 3' got: '%s'", actual)
		}
	})

	t.Run("register for boolean receiver", func(t *testing.T) {
		err := RegisterBoolFunc("negate", func(b bool, args ...any) bool {
			return !b
		})
		if err != nil {
			t.Fatalf("error registering function: %s", err)
		}

		actual, err := EvaluateString("{{ true.negate() }}", nil)
		if err != nil {
			t.Fatalf("error evaluating template: %s", err)
		}

		if actual != "0" {
			t.Fatalf("wrong result. expect: '0' got '%s'", actual)
		}
	})

	t.Run("register for string receiver", func(t *testing.T) {
		err := RegisterStrFunc("concat", func(s string, args ...any) string {
			arg1Value := args[0].(string)
			arg2Value := args[1].(string)

			return s + arg1Value + arg2Value
		})

		if err != nil {
			t.Fatalf("error registering function: %s", err)
		}

		actual, err := EvaluateString("{{ 'anna'.concat(' ', 'cho') }}", nil)
		if err != nil {
			t.Fatalf("error evaluating template: %s", err)
		}

		if actual != "anna cho" {
			t.Fatalf("wrong result. expect: 'anna cho' got: '%s'", actual)
		}
	})

	t.Run("registering already registered function", func(t *testing.T) {
		err := RegisterStrFunc("len", func(s string, args ...any) string {
			return "some output"
		})

		if err != nil {
			t.Fatalf("error registering function: %s", err)
		}

		// Registering the same function again should return an error
		err = RegisterStrFunc("len", func(s string, args ...any) string {
			return "some output"
		})

		if err == nil {
			t.Fatalf("expect error but got none")
		}

		expect := fail.New(0, "", "API", fail.ErrFuncAlreadyDefined, "len", "strings")

		if err.Error() != expect.Error().Error() {
			t.Fatalf("wrong error message. expect: %q got: %q", expect, err)
		}
	})

	t.Run("redefining built-in function not working", func(t *testing.T) {
		err := RegisterStrFunc("trim", func(s string, args ...any) string {
			return "some output"
		})

		if err != nil {
			t.Fatalf("error registering function: %s", err)
		}

		actual, err := EvaluateString("{{ ' anna '.trim() }}", nil)
		if err != nil {
			t.Fatalf("error registering function: %s", err)
		}

		// the output should be the same as the built-in function even though we redefined it
		if actual != "anna" {
			t.Fatalf("wrong output. expect: 'anna' got: '%s'", actual)
		}
	})
}
