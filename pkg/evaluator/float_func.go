package evaluator

import (
	"math"

	"github.com/textwire/textwire/v3/pkg/utils"
	"github.com/textwire/textwire/v3/pkg/value"
)

// floatIntFunc returns the integer part of the given float
func floatIntFunc(receiver value.Literal, _ ...value.Literal) (value.Literal, error) {
	val := receiver.(*value.Float).Val
	return &value.Int{Val: int64(val)}, nil
}

// floatStrFunc converts a float to a string and returns it
func floatStrFunc(receiver value.Literal, _ ...value.Literal) (value.Literal, error) {
	val := receiver.(*value.Float).Val
	return &value.Str{Val: utils.FloatToStr(val)}, nil
}

// floatAbsFunc returns the absolute value of an float
func floatAbsFunc(receiver value.Literal, _ ...value.Literal) (value.Literal, error) {
	val := receiver.(*value.Float).Val
	if val < 0 {
		return &value.Float{Val: -val}, nil
	}

	return receiver, nil
}

// floatCeilFunc returns the rounded up value of a float to the nearest integer
func floatCeilFunc(receiver value.Literal, _ ...value.Literal) (value.Literal, error) {
	val := receiver.(*value.Float).Val
	return &value.Int{Val: int64(math.Ceil(val))}, nil
}

// floatFloorFunc returns the rounded down value of a float to the nearest integer
func floatFloorFunc(receiver value.Literal, _ ...value.Literal) (value.Literal, error) {
	val := receiver.(*value.Float).Val
	return &value.Int{Val: int64(math.Floor(val))}, nil
}

// floatRoundFunc returns the rounded value of a float to the nearest integer
func floatRoundFunc(receiver value.Literal, _ ...value.Literal) (value.Literal, error) {
	val := receiver.(*value.Float).Val
	return &value.Int{Val: int64(math.Round(val))}, nil
}
