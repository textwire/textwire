package evaluator

import (
	"strconv"

	"github.com/textwire/textwire/v2/object"
)

// intFloatFunc converts an integer to a float and returns it
func intFloatFunc(receiver object.Object, args ...object.Object) (object.Object, error) {
	floatVal := receiver.(*object.Int).Value
	return &object.Float{Value: float64(floatVal)}, nil
}

// intAbsFunc returns the absolute value of an integer
func intAbsFunc(receiver object.Object, args ...object.Object) (object.Object, error) {
	intVal := receiver.(*object.Int).Value

	if intVal < 0 {
		return &object.Int{Value: -intVal}, nil
	}

	return receiver, nil
}

// intStrFunc converts an integer to a string and returns it
func intStrFunc(receiver object.Object, args ...object.Object) (object.Object, error) {
	intVal := receiver.(*object.Int).Value
	return &object.Str{Value: strconv.FormatInt(intVal, 10)}, nil
}
