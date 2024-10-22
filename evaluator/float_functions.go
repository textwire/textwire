package evaluator

import (
	"math"

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

// floatCeilFunc returns the rounded up value of a float to the nearest integer
func floatCeilFunc(receiver object.Object, args ...object.Object) (object.Object, error) {
	floatVal := receiver.(*object.Float).Value
	return &object.Int{Value: int64(math.Ceil(floatVal))}, nil
}

// floatFloorFunc returns the rounded down value of a float to the nearest integer
func floatFloorFunc(receiver object.Object, args ...object.Object) (object.Object, error) {
	floatVal := receiver.(*object.Float).Value
	return &object.Int{Value: int64(math.Floor(floatVal))}, nil
}

// floatRoundFunc returns the rounded value of a float to the nearest integer
func floatRoundFunc(receiver object.Object, args ...object.Object) (object.Object, error) {
	floatVal := receiver.(*object.Float).Value
	return &object.Int{Value: int64(math.Round(floatVal))}, nil
}
