package evaluator

import (
	"math"

	"github.com/textwire/textwire/v2/ctx"
	"github.com/textwire/textwire/v2/object"
	"github.com/textwire/textwire/v2/utils"
)

// floatIntFunc returns the integer part of the given float
func floatIntFunc(_ *ctx.EvalCtx, receiver object.Object, _ ...object.Object) (object.Object, error) {
	val := receiver.(*object.Float).Value
	return &object.Int{Value: int64(val)}, nil
}

// floatStrFunc converts a float to a string and returns it
func floatStrFunc(_ *ctx.EvalCtx, receiver object.Object, _ ...object.Object) (object.Object, error) {
	val := receiver.(*object.Float).Value
	return &object.Str{Value: utils.FloatToStr(val)}, nil
}

// floatAbsFunc returns the absolute value of an float
func floatAbsFunc(_ *ctx.EvalCtx, receiver object.Object, _ ...object.Object) (object.Object, error) {
	val := receiver.(*object.Float).Value

	if val < 0 {
		return &object.Float{Value: -val}, nil
	}

	return receiver, nil
}

// floatCeilFunc returns the rounded up value of a float to the nearest integer
func floatCeilFunc(_ *ctx.EvalCtx, receiver object.Object, _ ...object.Object) (object.Object, error) {
	val := receiver.(*object.Float).Value
	return &object.Int{Value: int64(math.Ceil(val))}, nil
}

// floatFloorFunc returns the rounded down value of a float to the nearest integer
func floatFloorFunc(_ *ctx.EvalCtx, receiver object.Object, _ ...object.Object) (object.Object, error) {
	val := receiver.(*object.Float).Value
	return &object.Int{Value: int64(math.Floor(val))}, nil
}

// floatRoundFunc returns the rounded value of a float to the nearest integer
func floatRoundFunc(_ *ctx.EvalCtx, receiver object.Object, _ ...object.Object) (object.Object, error) {
	val := receiver.(*object.Float).Value
	return &object.Int{Value: int64(math.Round(val))}, nil
}
