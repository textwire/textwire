package evaluator

import "github.com/textwire/textwire/v2/object"

// floatIntFunc returns the integer part of the given float
func floatIntFunc(receiver object.Object, args ...object.Object) (object.Object, error) {
	floatVal := receiver.(*object.Float).Value
	return &object.Int{Value: int64(floatVal)}, nil
}
