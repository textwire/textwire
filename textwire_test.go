package textwire

import (
	"io"
	"os"
	"testing"

	fail "github.com/textwire/textwire/v2/fail"
	object "github.com/textwire/textwire/v2/object"
)

func readFile(fileName string) (string, error) {
	file, err := os.Open(fileName)

	if err != nil {
		return "", err
	}

	defer file.Close()

	bytes, err := io.ReadAll(file)

	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func TestEvaluateString(t *testing.T) {
	tests := []struct {
		inp    string
		expect string
		data   map[string]interface{}
	}{
		{"{{ 1 + 2 }}", "3", nil},
		{"{{ n1 * n2 }}", "2", map[string]interface{}{"n1": 1, "n2": 2}},
		{"{{ user.name.firstName }}", "Ann", map[string]interface{}{"user": struct {
			Name struct{ FirstName string }
			Age  int
		}{Name: struct{ FirstName string }{"Ann"}, Age: 20}}},
	}

	for _, tt := range tests {
		actual, err := EvaluateString(tt.inp, tt.data)

		if err != nil {
			t.Errorf("error evaluating template: %s", err)
		}

		if actual != tt.expect {
			t.Errorf("wrong result. EXPECTED:\n\"%s\"\nGOT:\n\"%s\"",
				tt.expect, actual)
		}
	}
}

func TestErrorHandlingEvaluatingString(t *testing.T) {
	tests := []struct {
		inp  string
		err  *fail.Error
		data map[string]interface{}
	}{
		{`{{ 1 + "a" }}`, fail.New(1, "", "evaluator", fail.ErrTypeMismatch, object.INT_OBJ, "+", object.STR_OBJ), nil},
		{`@use("sometemplate")`, fail.New(1, "", "evaluator", fail.ErrUseStmtMustHaveProgram), nil},
		{`{{ loop = "test" }}`, fail.New(1, "", "evaluator", fail.ErrLoopVariableIsReserved), nil},
		{`{{ loop }}`, fail.New(0, "", "evaluator", fail.ErrLoopVariableIsReserved), map[string]interface{}{"loop": "test"}},
		{`{{ n = 1; n = "test" }}`, fail.New(1, "", "evaluator", fail.ErrVariableTypeMismatch, "n", object.INT_OBJ, object.STR_OBJ), nil},
		{`{{ obj = {}; obj.name }}`, fail.New(1, "", "evaluator", fail.ErrPropertyNotFound, "name", object.OBJ_OBJ), nil},
	}

	for _, tt := range tests {
		_, err := EvaluateString(tt.inp, tt.data)

		if err == nil {
			t.Errorf("expected error but got none")
			return
		}

		if err.Error() != tt.err.String() {
			t.Errorf("wrong error message. EXPECTED:\n%q\nGOT:\n%q",
				tt.err, err)
		}
	}
}

func TestEvaluateFile(t *testing.T) {
	absPath, fileErr := getFullPath("testdata/good/before/2.with-inserts", true)

	if fileErr != nil {
		t.Errorf("error getting full path: %s", fileErr)
		return
	}

	_, err := EvaluateFile(absPath, nil)

	if err != nil {
		t.Errorf("error evaluating file:\n%s", err)
	}
}

func TestEvaluateStringWithCustomFunction(t *testing.T) {
	t.Run("integer receiver", func(tt *testing.T) {
		RegisterIntFunc("double", func(num *object.Int, args ...object.Object) object.Object {
			return &object.Int{Value: num.Value * 2}
		})

		actual, err := EvaluateString("{{ 3.double() }}", nil)

		if err != nil {
			t.Errorf("error evaluating template: %s", err)
		}

		if actual != "6" {
			t.Errorf("wrong result. EXPECTED:\n\"%s\"\nGOT:\n\"%s\"", "6", actual)
		}
	})

	t.Run("float receiver", func(tt *testing.T) {
		// TODO: implement
	})

	t.Run("array receiver", func(tt *testing.T) {
		// TODO: implement
	})

	t.Run("string receiver", func(tt *testing.T) {
		// TODO: implement
	})

	t.Run("boolean receiver", func(tt *testing.T) {
		// TODO: implement
	})

	t.Run("with 2 arguments", func(tt *testing.T) {
		RegisterStrFunc("concat", func(s *object.Str, args ...object.Object) object.Object {
			arg1Value := args[0].(*object.Str).Value
			arg2Value := args[1].(*object.Str).Value

			return &object.Str{Value: s.Value + arg1Value + arg2Value}
		})

		actual, err := EvaluateString("{{ 'anna'.concat(' ', 'smith') }}", nil)

		if err != nil {
			t.Errorf("error evaluating template: %s", err)
		}

		if actual != "anna smith" {
			t.Errorf("wrong result. EXPECTED:\n\"%s\"\nGOT:\n\"%s\"", "anna smith", actual)
		}
	})
}
