package evaluator

import (
	"errors"
	"fmt"
	"html"
	"strings"
	"unicode/utf8"

	"github.com/textwire/textwire/v3/pkg/fail"
	"github.com/textwire/textwire/v3/pkg/object"
)

const defaultCharTrim = "\t \n\r"

// strLenFunc returns the length of the given string
func strLenFunc(receiver object.Object, _ ...object.Object) (object.Object, error) {
	val := receiver.(*object.String).Val
	chars := []rune(val)
	return &object.Integer{Val: int64(len(chars))}, nil
}

// strSplitFunc returns a list of strings split by the given separator
func strSplitFunc(receiver object.Object, args ...object.Object) (object.Object, error) {
	separator := " "

	if len(args) > 0 {
		str, ok := args[0].(*object.String)
		if !ok {
			msg := fmt.Sprintf(fail.ErrFuncFirstArgStr, object.STRING_OBJ, "split")
			return nil, errors.New(msg)
		}

		separator = str.Val
	}

	val := receiver.(*object.String).Val
	stringItems := strings.Split(val, separator)

	var elems []object.Object

	for _, val := range stringItems {
		elems = append(elems, &object.String{Val: val})
	}

	return &object.Array{Elements: elems}, nil
}

// strRawFunc prevents escaping HTML tags in a string
func strRawFunc(receiver object.Object, _ ...object.Object) (object.Object, error) {
	val := receiver.(*object.String).Val
	return &object.String{Val: html.UnescapeString(val)}, nil
}

// strTrimFunc returns a string with leading and trailing whitespace removed
func strTrimFunc(receiver object.Object, args ...object.Object) (object.Object, error) {
	chars := defaultCharTrim

	if len(args) > 0 {
		str, ok := args[0].(*object.String)
		if !ok {
			msg := fmt.Sprintf(fail.ErrFuncFirstArgStr, object.STRING_OBJ, "trim")
			return nil, errors.New(msg)
		}

		chars = str.Val
	}

	val := receiver.(*object.String).Val

	return &object.String{Val: strings.Trim(val, chars)}, nil
}

// strUpperFunc returns a string with all characters in uppercase
func strUpperFunc(receiver object.Object, _ ...object.Object) (object.Object, error) {
	val := receiver.(*object.String).Val
	return &object.String{Val: strings.ToUpper(val)}, nil
}

// strLowerFunc returns a string with all characters in lowercase
func strLowerFunc(receiver object.Object, _ ...object.Object) (object.Object, error) {
	str := receiver.(*object.String)
	return &object.String{Val: strings.ToLower(str.Val)}, nil
}

// strCapitalizeFunc returns a string with the first character capitalized
func strCapitalizeFunc(receiver object.Object, _ ...object.Object) (object.Object, error) {
	val := receiver.(*object.String).Val
	if len(val) == 0 {
		return &object.String{Val: ""}, nil
	}

	newVal := strings.ToUpper(val[:1]) + val[1:]

	return &object.String{Val: newVal}, nil
}

// strReverseFunc returns a string with the characters reversed
func strReverseFunc(receiver object.Object, _ ...object.Object) (object.Object, error) {
	val := receiver.(*object.String).Val

	runes := []rune(val)
	n := len(runes)

	// Reverse the slice of runes
	for i := 0; i < n/2; i++ {
		runes[i], runes[n-1-i] = runes[n-1-i], runes[i]
	}

	return &object.String{Val: string(runes)}, nil
}

// strContainsFunc returns true if the string contains the given substring, false otherwise
func strContainsFunc(receiver object.Object, args ...object.Object) (object.Object, error) {
	if len(args) == 0 {
		msg := fmt.Sprintf(fail.ErrFuncMissingArg, object.STRING_OBJ, "contains")
		return nil, errors.New(msg)
	}

	firstArg, ok := args[0].(*object.String)
	if !ok {
		msg := fmt.Sprintf(fail.ErrFuncFirstArgStr, object.STRING_OBJ, "contains")
		return nil, errors.New(msg)
	}

	val := receiver.(*object.String).Val
	substr := firstArg.Val

	return nativeBoolToBoolObj(strings.Contains(val, substr)), nil
}

// strTruncateFunc truncates a string to a specified length and appends an ellipsis.
func strTruncateFunc(receiver object.Object, args ...object.Object) (object.Object, error) {
	// Validate that at least the limit argument is provided
	if len(args) == 0 {
		msg := fmt.Sprintf(fail.ErrFuncMissingArg, object.STRING_OBJ, "truncate")
		return nil, errors.New(msg)
	}

	// Validate that the first argument is an integer (the limit)
	firstArg, ok := args[0].(*object.Integer)
	if !ok {
		msg := fmt.Sprintf(fail.ErrFuncFirstArgInt, object.STRING_OBJ, "truncate")
		return nil, errors.New(msg)
	}

	val := receiver.(*object.String).Val

	// Handle negative limits by treating them as 0
	// This prevents slice bounds errors when limit < 0
	limit := max(int(firstArg.Val), 0)

	// If the string is already shorter than or equal to the limit, return it unchanged
	// This ensures we don't truncate strings that don't need truncation
	if limit >= utf8.RuneCountInString(val) {
		return &object.String{Val: val}, nil
	}

	ellipsis := "..."

	// If a custom suffix is provided as the second argument, use it instead
	if len(args) > 1 {
		secondArg, ok := args[1].(*object.String)
		if ok {
			ellipsis = secondArg.Val
		} else {
			msg := fmt.Sprintf(fail.ErrFuncSecondArgStr, object.STRING_OBJ, "truncate")
			return nil, errors.New(msg)
		}
	}

	// Truncate the string at the limit and append the suffix
	return &object.String{Val: val[:limit] + ellipsis}, nil
}

// strDecimalFunc returns a string formatted as a decimal number.
// For integers: appends decimal places (e.g., "100" → "100.00")
// For floats: ensures minimum decimal places (e.g., "-0.5" → "-0.50")
// Non-numeric strings are returned unchanged.
func strDecimalFunc(receiver object.Object, args ...object.Object) (object.Object, error) {
	val := receiver.(*object.String).Val
	separator, decimals, err := getDecimalConfig(object.STRING_OBJ, args...)
	if err != nil {
		return nil, err
	}

	switch {
	case strIsInt(val):
		return &object.String{Val: formatIntDecimals(val, separator, decimals)}, nil
	case isValidFloat(val):
		return &object.String{Val: formatFloatDecimals(val, separator, decimals)}, nil
	default:
		return &object.String{Val: val}, nil
	}
}

// strAtFunc returns the character at the given index in the string
func strAtFunc(receiver object.Object, args ...object.Object) (object.Object, error) {
	index := 0

	if len(args) != 0 {
		firstArg, ok := args[0].(*object.Integer)
		if !ok {
			msg := fmt.Sprintf(fail.ErrFuncFirstArgInt, object.STRING_OBJ, "at")
			return nil, errors.New(msg)
		}

		index = int(firstArg.Val)
	}

	val := receiver.(*object.String).Val

	chars := []rune(val)
	if len(chars) == 0 {
		return &object.Nil{}, nil
	}

	if index < 0 {
		index = len(chars) + index
	}

	if index >= len(chars) {
		return &object.Nil{}, nil
	}

	return &object.String{Val: string(chars[index])}, nil
}

// strFirstFunc returns the first character in the string
func strFirstFunc(receiver object.Object, _ ...object.Object) (object.Object, error) {
	return strAtFunc(receiver, &object.Integer{Val: 0})
}

// strLastFunc returns the last character in the string
func strLastFunc(receiver object.Object, _ ...object.Object) (object.Object, error) {
	return strAtFunc(receiver, &object.Integer{Val: -1})
}

// strTrimRightFunc returns a string with trailing whitespace removed from the right
func strTrimRightFunc(receiver object.Object, args ...object.Object) (object.Object, error) {
	chars := defaultCharTrim
	if len(args) > 0 {
		str, ok := args[0].(*object.String)
		if !ok {
			msg := fmt.Sprintf(fail.ErrFuncFirstArgStr, object.STRING_OBJ, "trimRight")
			return nil, errors.New(msg)
		}

		chars = str.Val
	}

	val := receiver.(*object.String).Val

	return &object.String{Val: strings.TrimRight(val, chars)}, nil
}

// strTrimLeftFunc returns a string with trailing whitespace removed from the left
func strTrimLeftFunc(receiver object.Object, args ...object.Object) (object.Object, error) {
	chars := defaultCharTrim
	if len(args) > 0 {
		str, ok := args[0].(*object.String)

		if !ok {
			msg := fmt.Sprintf(fail.ErrFuncFirstArgStr, object.STRING_OBJ, "trimLeft")
			return nil, errors.New(msg)
		}

		chars = str.Val
	}

	val := receiver.(*object.String).Val

	return &object.String{Val: strings.TrimLeft(val, chars)}, nil
}

// strRepeatFunc returns a string repeated n times
func strRepeatFunc(receiver object.Object, args ...object.Object) (object.Object, error) {
	if len(args) == 0 {
		msg := fmt.Sprintf(fail.ErrFuncMissingArg, object.STRING_OBJ, "repeat")
		return nil, errors.New(msg)
	}

	firstArg, ok := args[0].(*object.Integer)
	if !ok {
		msg := fmt.Sprintf(fail.ErrFuncFirstArgInt, object.STRING_OBJ, "repeat")
		return nil, errors.New(msg)
	}

	val := receiver.(*object.String).Val
	count := int(firstArg.Val)

	if count < 0 {
		return &object.String{Val: ""}, nil
	}

	repeated := strings.Repeat(val, count)

	return &object.String{Val: repeated}, nil
}

// strFormatFunc embeds values into a string. Similar to sprintf in C.
func strFormatFunc(receiver object.Object, args ...object.Object) (object.Object, error) {
	if len(args) == 0 {
		msg := fmt.Sprintf(fail.ErrFuncMissingArg, object.STRING_OBJ, "format")
		return nil, errors.New(msg)
	}

	str := receiver.(*object.String).Val

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

	return &object.String{Val: out.String()}, nil
}
