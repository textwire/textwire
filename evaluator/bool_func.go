package evaluator

import (
	"github.com/textwire/textwire/v2/object"
)

// boolBinaryFunc returns an integer 1 if the receiver is true, 0 otherwise
func boolBinaryFunc(receiver object.Object, _ ...object.Object) (object.Object, error) {
	if isTruthy(receiver) {
		return &object.Int{Value: 1}, nil
	}

	return &object.Int{Value: 0}, nil
}
