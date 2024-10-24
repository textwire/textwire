package evaluator

import (
	"strconv"

	"github.com/textwire/textwire/v2/object"
)

// intFloatFunc converts an integer to a float and returns it
func intFloatFunc(receiver object.Object, _ ...object.Object) (object.Object, error) {
	val := receiver.(*object.Int).Value
	return &object.Float{Value: float64(val)}, nil
}

// intAbsFunc returns the absolute value of an integer
func intAbsFunc(receiver object.Object, _ ...object.Object) (object.Object, error) {
	val := receiver.(*object.Int).Value

	if val < 0 {
		return &object.Int{Value: -val}, nil
	}

	return receiver, nil
}

// intStrFunc converts an integer to a string and returns it
func intStrFunc(receiver object.Object, _ ...object.Object) (object.Object, error) {
	val := receiver.(*object.Int).Value
	return &object.Str{Value: strconv.FormatInt(val, 10)}, nil
}
