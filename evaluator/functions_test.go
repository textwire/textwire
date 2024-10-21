package evaluator

import (
	"testing"

	"github.com/textwire/textwire/v2/fail"
	"github.com/textwire/textwire/v2/object"
)

func TestFunctionGivesError(t *testing.T) {
	tests := []struct {
		inp         string
		expectedErr string
		funcName    string
	}{
		{`{{ [1, 2].slice() }}`, fail.ErrFuncRequiresOneArg, "slice"},
	}

	for _, tc := range tests {
		evaluated := testEval(tc.inp)
		errObj, ok := evaluated.(*object.Error)

		if !ok {
			t.Fatalf("evaluation failed: %s", errObj.String())
		}

		if evaluated.Type() != object.ERR_OBJ {
			t.Fatalf("expected object.ERR_OBJ, got=%T", evaluated)
		}

		failErr := fail.New(1, "/path/to/file", "evaluator", tc.expectedErr, tc.funcName)

		if errObj.String() != failErr.String() {
			t.Fatalf("expected error message=%q, got=%q", failErr.String(), errObj.String())
		}
	}
}
