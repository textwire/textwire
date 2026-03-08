package evaluator

import (
	"errors"
	"fmt"

	"github.com/textwire/textwire/v3/pkg/fail"
	"github.com/textwire/textwire/v3/pkg/value"
)

// boolBinaryFunc returns an integer 1 if the receiver is true, 0 otherwise
func boolBinaryFunc(receiver value.Value, _ ...value.Value) (value.Value, error) {
	if isTruthy(receiver) {
		return &value.Int{Val: 1}, nil
	}

	return &value.Int{Val: 0}, nil
}

// boolThenFunc returns the first argument if the receiver is true, the second argument or nil otherwise
func boolThenFunc(receiver value.Value, args ...value.Value) (value.Value, error) {
	if len(args) == 0 {
		msg := fmt.Sprintf(fail.ErrFuncMissingArg, value.BOOL_OBJ, "then")
		return nil, errors.New(msg)
	}

	if isTruthy(receiver) {
		return args[0], nil
	}

	if len(args) == 1 {
		return &value.Nil{}, nil
	}

	return args[1], nil
}
