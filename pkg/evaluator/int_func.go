package evaluator

import (
	"strconv"

	"github.com/textwire/textwire/v3/pkg/value"
)

// intFloatFunc converts an integer to a float and returns it
func intFloatFunc(receiver value.Literal, _ ...value.Literal) (value.Literal, error) {
	val := receiver.(*value.Int).Val
	return &value.Float{Val: float64(val)}, nil
}

// intAbsFunc returns the absolute value of an integer
func intAbsFunc(receiver value.Literal, _ ...value.Literal) (value.Literal, error) {
	val := receiver.(*value.Int).Val
	if val < 0 {
		return &value.Int{Val: -val}, nil
	}

	return receiver, nil
}

// intStrFunc converts an integer to a string and returns it
func intStrFunc(receiver value.Literal, _ ...value.Literal) (value.Literal, error) {
	val := receiver.(*value.Int).Val
	return &value.Str{Val: strconv.FormatInt(val, 10)}, nil
}

// intLenFunc returns the number of digits in an integer
func intLenFunc(receiver value.Literal, _ ...value.Literal) (value.Literal, error) {
	val := receiver.(*value.Int).Val
	valStr := strconv.FormatInt(val, 10)
	if val < 0 {
		return &value.Int{Val: int64(len(valStr) - 1)}, nil
	}

	return &value.Int{Val: int64(len(valStr))}, nil
}

// intDecimalFunc returns a string formatted as a decimal number.
// Converts integer to string and appends decimal places (e.g., 100 → "100.00")
func intDecimalFunc(receiver value.Literal, args ...value.Literal) (value.Literal, error) {
	val := receiver.(*value.Int).String()
	separator, decimals, err := getDecimalConfig(value.INT_VAL, args...)
	if err != nil {
		return nil, err
	}

	return &value.Str{Val: formatIntDecimals(val, separator, decimals)}, nil
}
