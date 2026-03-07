package evaluator

import (
	"strconv"

	"github.com/textwire/textwire/v3/pkg/object"
)

// intFloatFunc converts an integer to a float and returns it
func intFloatFunc(receiver object.Object, _ ...object.Object) (object.Object, error) {
	val := receiver.(*object.Integer).Val
	return &object.Float{Val: float64(val)}, nil
}

// intAbsFunc returns the absolute value of an integer
func intAbsFunc(receiver object.Object, _ ...object.Object) (object.Object, error) {
	val := receiver.(*object.Integer).Val
	if val < 0 {
		return &object.Integer{Val: -val}, nil
	}

	return receiver, nil
}

// intStrFunc converts an integer to a string and returns it
func intStrFunc(receiver object.Object, _ ...object.Object) (object.Object, error) {
	val := receiver.(*object.Integer).Val
	return &object.String{Val: strconv.FormatInt(val, 10)}, nil
}

// intLenFunc returns the number of digits in an integer
func intLenFunc(receiver object.Object, _ ...object.Object) (object.Object, error) {
	val := receiver.(*object.Integer).Val
	valStr := strconv.FormatInt(val, 10)
	if val < 0 {
		return &object.Integer{Val: int64(len(valStr) - 1)}, nil
	}

	return &object.Integer{Val: int64(len(valStr))}, nil
}

// intDecimalFunc returns a string formatted as a decimal number.
// Converts integer to string and appends decimal places (e.g., 100 → "100.00")
func intDecimalFunc(receiver object.Object, args ...object.Object) (object.Object, error) {
	val := receiver.(*object.Integer).String()
	separator, decimals, err := getDecimalConfig(object.INTEGER_OBJ, args...)
	if err != nil {
		return nil, err
	}

	return &object.String{Val: formatIntDecimals(val, separator, decimals)}, nil
}
