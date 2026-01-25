package textwire

import (
	"io"
	"log"
	"os"
	"testing"

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
		inp    string
		expect string
		data   map[string]any
	}{
		{"{{ 1 + 2 }}", "3", nil},
		{"{{ n1 * n2 }}", "2", map[string]any{"n1": 1, "n2": 2}},
		{"{{ user.name.firstName }}", "Ann", map[string]any{"user": struct {
			Name struct{ FirstName string }
			Age  int
		}{Name: struct{ FirstName string }{"Ann"}, Age: 20}}},
	}

	for _, tc := range cases {
		actual, err := EvaluateString(tc.inp, tc.data)
		if err != nil {
			t.Errorf("error evaluating template: %s", err)
		}

		if actual != tc.expect {
			t.Errorf("wrong result. EXPECTED:\n\"%s\"\nGOT:\n\"%s\"",
				tc.expect, actual)
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
		{`@use("sometemplate")`, fail.New(1, "", "evaluator", fail.ErrUseStmtMustHaveProgram), nil},
		{`{{ loop = "test" }}`, fail.New(1, "", "evaluator", fail.ErrReservedVariables), nil},
		{`{{ global = "test" }}`, fail.New(1, "", "evaluator", fail.ErrReservedVariables), nil},
		{`{{ loop }}`, fail.New(0, "", "evaluator", fail.ErrReservedVariables), map[string]any{"loop": "test"}},
		{`{{ n = 1; n = "test" }}`, fail.New(1, "", "evaluator", fail.ErrVariableTypeMismatch, "n", object.INT_OBJ, object.STR_OBJ), nil},
		{`{{ obj = {}; obj.name }}`, fail.New(1, "", "evaluator", fail.ErrPropertyNotFound, "name", object.OBJ_OBJ), nil},
		{`{{ {}.test }}`, fail.New(1, "", "evaluator", fail.ErrPropertyNotFound, "test", object.OBJ_OBJ), nil},
		{`{{ 5.somefunction() }}`, fail.New(1, "", "evaluator", fail.ErrNoFuncForThisType, "somefunction", object.INT_OBJ), nil},
		{`{{ 3 / 0 }}`, fail.New(1, "", "evaluator", fail.ErrDivisionByZero), nil},
		{`{{ 1 ~ 8 }}`, fail.New(1, "", "parser", fail.ErrIllegalToken, "~"), nil},
	}

	for _, tc := range cases {
		_, err := EvaluateString(tc.inp, tc.data)
		if err == nil {
			t.Errorf("expected error but got none")
			return
		}

		if err.Error() != tc.err.String() {
			t.Errorf("wrong error message. EXPECTED:\n%q\nGOT:\n%q",
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

	expected, err := readFile("textwire/testdata/good/expected/" + filename + ".html")
	if err != nil {
		t.Errorf("error reading expected file: %s", err)
		return
	}

	if out != expected {
		t.Errorf("wrong output. EXPECTED:\n%s\nGOT:\n%s", expected, out)
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
			t.Errorf("wrong result. EXPECTED: '6' GOT: '%s'", actual)
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
			t.Fatalf("wrong result. EXPECTED: '7.0' GOT: '%s'", actual)
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
			t.Fatalf("wrong result. EXPECTED: '1, 2, 3' GOT: '%s'", actual)
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
			t.Fatalf("wrong result. EXPECTED: '0' GOT '%s'", actual)
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
			t.Fatalf("wrong result. EXPECTED: 'anna cho' GOT: '%s'", actual)
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
			t.Fatalf("expected error but got none")
		}

		expect := fail.New(0, "", "API", fail.ErrFuncAlreadyDefined, "len", "strings")

		if err.Error() != expect.Error().Error() {
			t.Fatalf("wrong error message. EXPECTED: %q GOT: %q", expect, err)
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
			t.Fatalf("wrong output. EXPECTED: 'anna' GOT: '%s'", actual)
		}
	})
}
