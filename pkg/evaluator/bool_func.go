package evaluator

import (
	"errors"
	"fmt"

	"github.com/textwire/textwire/v3/pkg/fail"
	"github.com/textwire/textwire/v3/pkg/object"
)

// boolBinaryFunc returns an integer 1 if the receiver is true, 0 otherwise
func boolBinaryFunc(receiver object.Object, _ ...object.Object) (object.Object, error) {
	if isTruthy(receiver) {
		return &object.Int{Value: 1}, nil
	}

	return &object.Int{Value: 0}, nil
}

// boolThenFunc returns the first argument if the receiver is true, the second argument or nil otherwise
func boolThenFunc(receiver object.Object, args ...object.Object) (object.Object, error) {
	if len(args) == 0 {
		msg := fmt.Sprintf(fail.ErrFuncMissingArg, object.BOOL_OBJ, "then")
		return nil, errors.New(msg)
	}

	if isTruthy(receiver) {
		return args[0], nil
	}

	if len(args) == 1 {
		return &object.Nil{}, nil
	}

	return args[1], nil
}
