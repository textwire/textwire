package evaluator

import (
	"errors"
	"fmt"
	"html"
	"strings"
	"unicode/utf8"

	"github.com/textwire/textwire/v3/pkg/fail"
	"github.com/textwire/textwire/v3/pkg/value"
)

const defaultCharTrim = "\t \n\r"

// strLenFunc returns the length of the given string
func strLenFunc(receiver value.Value, _ ...value.Value) (value.Value, error) {
	val := receiver.(*value.Str).Val
	chars := []rune(val)
	return &value.Int{Val: int64(len(chars))}, nil
}

// strSplitFunc returns a list of strings split by the given separator
func strSplitFunc(receiver value.Value, args ...value.Value) (value.Value, error) {
	separator := " "

	if len(args) > 0 {
		str, ok := args[0].(*value.Str)
		if !ok {
			msg := fmt.Sprintf(fail.ErrFuncFirstArgStr, value.STR_VAL, "split")
			return nil, errors.New(msg)
		}

		separator = str.Val
	}

	val := receiver.(*value.Str).Val
	stringItems := strings.Split(val, separator)

	var elems []value.Value

	for _, val := range stringItems {
		elems = append(elems, &value.Str{Val: val})
	}

	return &value.Arr{Elements: elems}, nil
}

// strRawFunc prevents escaping HTML tags in a string
func strRawFunc(receiver value.Value, _ ...value.Value) (value.Value, error) {
	val := receiver.(*value.Str).Val
	return &value.Str{Val: html.UnescapeString(val)}, nil
}

// strTrimFunc returns a string with leading and trailing whitespace removed
func strTrimFunc(receiver value.Value, args ...value.Value) (value.Value, error) {
	chars := defaultCharTrim

	if len(args) > 0 {
		str, ok := args[0].(*value.Str)
		if !ok {
			msg := fmt.Sprintf(fail.ErrFuncFirstArgStr, value.STR_VAL, "trim")
			return nil, errors.New(msg)
		}

		chars = str.Val
	}

	val := receiver.(*value.Str).Val

	return &value.Str{Val: strings.Trim(val, chars)}, nil
}

// strUpperFunc returns a string with all characters in uppercase
func strUpperFunc(receiver value.Value, _ ...value.Value) (value.Value, error) {
	val := receiver.(*value.Str).Val
	return &value.Str{Val: strings.ToUpper(val)}, nil
}

// strLowerFunc returns a string with all characters in lowercase
func strLowerFunc(receiver value.Value, _ ...value.Value) (value.Value, error) {
	str := receiver.(*value.Str)
	return &value.Str{Val: strings.ToLower(str.Val)}, nil
}

// strCapitalizeFunc returns a string with the first character capitalized
func strCapitalizeFunc(receiver value.Value, _ ...value.Value) (value.Value, error) {
	val := receiver.(*value.Str).Val
	if len(val) == 0 {
		return &value.Str{Val: ""}, nil
	}

	newVal := strings.ToUpper(val[:1]) + val[1:]

	return &value.Str{Val: newVal}, nil
}

// strReverseFunc returns a string with the characters reversed
func strReverseFunc(receiver value.Value, _ ...value.Value) (value.Value, error) {
	val := receiver.(*value.Str).Val

	runes := []rune(val)
	n := len(runes)

	// Reverse the slice of runes
	for i := 0; i < n/2; i++ {
		runes[i], runes[n-1-i] = runes[n-1-i], runes[i]
	}

	return &value.Str{Val: string(runes)}, nil
}

// strContainsFunc returns true if the string contains the given substring, false otherwise
func strContainsFunc(receiver value.Value, args ...value.Value) (value.Value, error) {
	if len(args) == 0 {
		msg := fmt.Sprintf(fail.ErrFuncMissingArg, value.STR_VAL, "contains")
		return nil, errors.New(msg)
	}

	firstArg, ok := args[0].(*value.Str)
	if !ok {
		msg := fmt.Sprintf(fail.ErrFuncFirstArgStr, value.STR_VAL, "contains")
		return nil, errors.New(msg)
	}

	val := receiver.(*value.Str).Val
	substr := firstArg.Val

	return nativeBoolToBoolObj(strings.Contains(val, substr)), nil
}

// strTruncateFunc truncates a string to a specified length and appends an ellipsis.
func strTruncateFunc(receiver value.Value, args ...value.Value) (value.Value, error) {
	// Validate that at least the limit argument is provided
	if len(args) == 0 {
		msg := fmt.Sprintf(fail.ErrFuncMissingArg, value.STR_VAL, "truncate")
		return nil, errors.New(msg)
	}

	// Validate that the first argument is an integer (the limit)
	firstArg, ok := args[0].(*value.Int)
	if !ok {
		msg := fmt.Sprintf(fail.ErrFuncFirstArgInt, value.STR_VAL, "truncate")
		return nil, errors.New(msg)
	}

	val := receiver.(*value.Str).Val

	// Handle negative limits by treating them as 0
	// This prevents slice bounds errors when limit < 0
	limit := max(int(firstArg.Val), 0)

	// If the string is already shorter than or equal to the limit, return it unchanged
	// This ensures we don't truncate strings that don't need truncation
	if limit >= utf8.RuneCountInString(val) {
		return &value.Str{Val: val}, nil
	}

	ellipsis := "..."

	// If a custom suffix is provided as the second argument, use it instead
	if len(args) > 1 {
		secondArg, ok := args[1].(*value.Str)
		if ok {
			ellipsis = secondArg.Val
		} else {
			msg := fmt.Sprintf(fail.ErrFuncSecondArgStr, value.STR_VAL, "truncate")
			return nil, errors.New(msg)
		}
	}

	// Truncate the string at the limit and append the suffix
	return &value.Str{Val: val[:limit] + ellipsis}, nil
}

// strDecimalFunc returns a string formatted as a decimal number.
// For integers: appends decimal places (e.g., "100" → "100.00")
// For floats: ensures minimum decimal places (e.g., "-0.5" → "-0.50")
// Non-numeric strings are returned unchanged.
func strDecimalFunc(receiver value.Value, args ...value.Value) (value.Value, error) {
	val := receiver.(*value.Str).Val
	separator, decimals, err := getDecimalConfig(value.STR_VAL, args...)
	if err != nil {
		return nil, err
	}

	switch {
	case strIsInt(val):
		return &value.Str{Val: formatIntDecimals(val, separator, decimals)}, nil
	case isValidFloat(val):
		return &value.Str{Val: formatFloatDecimals(val, separator, decimals)}, nil
	default:
		return &value.Str{Val: val}, nil
	}
}

// strAtFunc returns the character at the given index in the string
func strAtFunc(receiver value.Value, args ...value.Value) (value.Value, error) {
	index := 0

	if len(args) != 0 {
		firstArg, ok := args[0].(*value.Int)
		if !ok {
			msg := fmt.Sprintf(fail.ErrFuncFirstArgInt, value.STR_VAL, "at")
			return nil, errors.New(msg)
		}

		index = int(firstArg.Val)
	}

	val := receiver.(*value.Str).Val

	chars := []rune(val)
	if len(chars) == 0 {
		return &value.Nil{}, nil
	}

	if index < 0 {
		index = len(chars) + index
	}

	if index >= len(chars) {
		return &value.Nil{}, nil
	}

	return &value.Str{Val: string(chars[index])}, nil
}

// strFirstFunc returns the first character in the string
func strFirstFunc(receiver value.Value, _ ...value.Value) (value.Value, error) {
	return strAtFunc(receiver, &value.Int{Val: 0})
}

// strLastFunc returns the last character in the string
func strLastFunc(receiver value.Value, _ ...value.Value) (value.Value, error) {
	return strAtFunc(receiver, &value.Int{Val: -1})
}

// strTrimRightFunc returns a string with trailing whitespace removed from the right
func strTrimRightFunc(receiver value.Value, args ...value.Value) (value.Value, error) {
	chars := defaultCharTrim
	if len(args) > 0 {
		str, ok := args[0].(*value.Str)
		if !ok {
			msg := fmt.Sprintf(fail.ErrFuncFirstArgStr, value.STR_VAL, "trimRight")
			return nil, errors.New(msg)
		}

		chars = str.Val
	}

	val := receiver.(*value.Str).Val

	return &value.Str{Val: strings.TrimRight(val, chars)}, nil
}

// strTrimLeftFunc returns a string with trailing whitespace removed from the left
func strTrimLeftFunc(receiver value.Value, args ...value.Value) (value.Value, error) {
	chars := defaultCharTrim
	if len(args) > 0 {
		str, ok := args[0].(*value.Str)

		if !ok {
			msg := fmt.Sprintf(fail.ErrFuncFirstArgStr, value.STR_VAL, "trimLeft")
			return nil, errors.New(msg)
		}

		chars = str.Val
	}

	val := receiver.(*value.Str).Val

	return &value.Str{Val: strings.TrimLeft(val, chars)}, nil
}

// strRepeatFunc returns a string repeated n times
func strRepeatFunc(receiver value.Value, args ...value.Value) (value.Value, error) {
	if len(args) == 0 {
		msg := fmt.Sprintf(fail.ErrFuncMissingArg, value.STR_VAL, "repeat")
		return nil, errors.New(msg)
	}

	firstArg, ok := args[0].(*value.Int)
	if !ok {
		msg := fmt.Sprintf(fail.ErrFuncFirstArgInt, value.STR_VAL, "repeat")
		return nil, errors.New(msg)
	}

	val := receiver.(*value.Str).Val
	count := int(firstArg.Val)

	if count < 0 {
		return &value.Str{Val: ""}, nil
	}

	repeated := strings.Repeat(val, count)

	return &value.Str{Val: repeated}, nil
}

// strFormatFunc embeds values into a string. Similar to sprintf in C.
func strFormatFunc(receiver value.Value, args ...value.Value) (value.Value, error) {
	if len(args) == 0 {
		msg := fmt.Sprintf(fail.ErrFuncMissingArg, value.STR_VAL, "format")
		return nil, errors.New(msg)
	}

	str := receiver.(*value.Str).Val

	var argIdx int
	var out strings.Builder
	out.Grow(len(args))

	for i := 0; i < len(str); i++ {
		isPlaceholder := str[i] == '%' && i+1 < len(str) && str[i+1] == 's'
		if !isPlaceholder {
			out.WriteByte(str[i])
			continue
		}

		if argIdx >= len(args) {
			out.WriteByte(str[i])
			continue
		}

		argVal := args[argIdx].String()
		out.WriteString(argVal)
		argIdx++
		i++
	}

	return &value.Str{Val: out.String()}, nil
}
