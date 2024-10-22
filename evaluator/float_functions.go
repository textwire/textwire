package evaluator

import (
	"github.com/textwire/textwire/v2/object"
	"github.com/textwire/textwire/v2/utils"
)

// floatIntFunc returns the integer part of the given float
func floatIntFunc(receiver object.Object, args ...object.Object) (object.Object, error) {
	floatVal := receiver.(*object.Float).Value
	return &object.Int{Value: int64(floatVal)}, nil
}

// floatStrFunc converts a float to a string and returns it
func floatStrFunc(receiver object.Object, args ...object.Object) (object.Object, error) {
	val := receiver.(*object.Float).Value
	return &object.Str{Value: utils.FloatToStr(val)}, nil
}

// floatAbsFunc returns the absolute value of an float
func floatAbsFunc(receiver object.Object, args ...object.Object) (object.Object, error) {
	floatVal := receiver.(*object.Float).Value

	if floatVal < 0 {
		return &object.Float{Value: -floatVal}, nil
	}

	return receiver, nil
}
