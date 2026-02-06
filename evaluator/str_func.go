package evaluator

import (
	"errors"
	"fmt"
	"html"
	"strings"
	"unicode/utf8"

	"github.com/textwire/textwire/v3/fail"
	"github.com/textwire/textwire/v3/object"
)

const defaultCharTrim = "\t \n\r"

// strLenFunc returns the length of the given string
func strLenFunc(receiver object.Object, _ ...object.Object) (object.Object, error) {
	val := receiver.(*object.Str).Value
	chars := []rune(val)
	return &object.Int{Value: int64(len(chars))}, nil
}

// strSplitFunc returns a list of strings split by the given separator
func strSplitFunc(receiver object.Object, args ...object.Object) (object.Object, error) {
	separator := " "

	if len(args) > 0 {
		str, ok := args[0].(*object.Str)
		if !ok {
			msg := fmt.Sprintf(fail.ErrFuncFirstArgStr, object.STR_OBJ, "split")
			return nil, errors.New(msg)
		}

		separator = str.Value
	}

	val := receiver.(*object.Str).Value
	stringItems := strings.Split(val, separator)

	var elems []object.Object

	for _, val := range stringItems {
		elems = append(elems, &object.Str{Value: val})
	}

	return &object.Array{Elements: elems}, nil
}

// strRawFunc prevents escaping HTML tags in a string
func strRawFunc(receiver object.Object, _ ...object.Object) (object.Object, error) {
	val := receiver.(*object.Str).Value
	return &object.Str{Value: html.UnescapeString(val)}, nil
}

// strTrimFunc returns a string with leading and trailing whitespace removed
func strTrimFunc(receiver object.Object, args ...object.Object) (object.Object, error) {
	chars := defaultCharTrim

	if len(args) > 0 {
		str, ok := args[0].(*object.Str)
		if !ok {
			msg := fmt.Sprintf(fail.ErrFuncFirstArgStr, object.STR_OBJ, "trim")
			return nil, errors.New(msg)
		}

		chars = str.Value
	}

	val := receiver.(*object.Str).Value

	return &object.Str{Value: strings.Trim(val, chars)}, nil
}

// strUpperFunc returns a string with all characters in uppercase
func strUpperFunc(receiver object.Object, _ ...object.Object) (object.Object, error) {
	val := receiver.(*object.Str).Value
	return &object.Str{Value: strings.ToUpper(val)}, nil
}

// strLowerFunc returns a string with all characters in lowercase
func strLowerFunc(receiver object.Object, _ ...object.Object) (object.Object, error) {
	str := receiver.(*object.Str)
	return &object.Str{Value: strings.ToLower(str.Value)}, nil
}

// strCapitalizeFunc returns a string with the first character capitalized
func strCapitalizeFunc(receiver object.Object, _ ...object.Object) (object.Object, error) {
	val := receiver.(*object.Str).Value
	if len(val) == 0 {
		return &object.Str{Value: ""}, nil
	}

	newVal := strings.ToUpper(val[:1]) + val[1:]

	return &object.Str{Value: newVal}, nil
}

// strReverseFunc returns a string with the characters reversed
func strReverseFunc(receiver object.Object, _ ...object.Object) (object.Object, error) {
	val := receiver.(*object.Str).Value

	runes := []rune(val)
	n := len(runes)

	// Reverse the slice of runes
	for i := 0; i < n/2; i++ {
		runes[i], runes[n-1-i] = runes[n-1-i], runes[i]
	}

	return &object.Str{Value: string(runes)}, nil
}

// strContainsFunc returns true if the string contains the given substring, false otherwise
func strContainsFunc(receiver object.Object, args ...object.Object) (object.Object, error) {
	if len(args) == 0 {
		msg := fmt.Sprintf(fail.ErrFuncMissingArg, object.STR_OBJ, "contains")
		return nil, errors.New(msg)
	}

	firstArg, ok := args[0].(*object.Str)
	if !ok {
		msg := fmt.Sprintf(fail.ErrFuncFirstArgStr, object.STR_OBJ, "contains")
		return nil, errors.New(msg)
	}

	val := receiver.(*object.Str).Value
	substr := firstArg.Value

	return nativeBoolToBooleanObject(strings.Contains(val, substr)), nil
}

// strTruncateFunc returns a string truncated to the given length
func strTruncateFunc(receiver object.Object, args ...object.Object) (object.Object, error) {
	if len(args) == 0 {
		msg := fmt.Sprintf(fail.ErrFuncMissingArg, object.STR_OBJ, "truncate")
		return nil, errors.New(msg)
	}

	firstArg, ok := args[0].(*object.Int)
	if !ok {
		msg := fmt.Sprintf(fail.ErrFuncFirstArgInt, object.STR_OBJ, "truncate")
		return nil, errors.New(msg)
	}

	val := receiver.(*object.Str).Value

	limit := int(firstArg.Value)
	if limit >= utf8.RuneCountInString(val) {
		return &object.Str{Value: val}, nil
	}

	ellipsis := "..."

	if len(args) > 1 {
		secondArg, ok := args[1].(*object.Str)
		if ok {
			ellipsis = secondArg.Value
		} else {
			msg := fmt.Sprintf(fail.ErrFuncSecondArgStr, object.STR_OBJ, "truncate")
			return nil, errors.New(msg)
		}
	}

	newVal := val[:firstArg.Value] + ellipsis

	return &object.Str{Value: newVal}, nil
}

// strDecimalFunc returns a string formatted as a decimal number
func strDecimalFunc(receiver object.Object, args ...object.Object) (object.Object, error) {
	return addDecimals(receiver, object.STR_OBJ, args...)
}

// strAtFunc returns the character at the given index in the string
func strAtFunc(receiver object.Object, args ...object.Object) (object.Object, error) {
	index := 0

	if len(args) != 0 {
		firstArg, ok := args[0].(*object.Int)
		if !ok {
			msg := fmt.Sprintf(fail.ErrFuncFirstArgInt, object.STR_OBJ, "at")
			return nil, errors.New(msg)
		}

		index = int(firstArg.Value)
	}

	val := receiver.(*object.Str).Value

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

	return &object.Str{Value: string(chars[index])}, nil
}

// strFirstFunc returns the first character in the string
func strFirstFunc(receiver object.Object, _ ...object.Object) (object.Object, error) {
	return strAtFunc(receiver, &object.Int{Value: 0})
}

// strLastFunc returns the last character in the string
func strLastFunc(receiver object.Object, _ ...object.Object) (object.Object, error) {
	return strAtFunc(receiver, &object.Int{Value: -1})
}

// strTrimRightFunc returns a string with trailing whitespace removed from the right
func strTrimRightFunc(receiver object.Object, args ...object.Object) (object.Object, error) {
	chars := defaultCharTrim
	if len(args) > 0 {
		str, ok := args[0].(*object.Str)
		if !ok {
			msg := fmt.Sprintf(fail.ErrFuncFirstArgStr, object.STR_OBJ, "trimRight")
			return nil, errors.New(msg)
		}

		chars = str.Value
	}

	val := receiver.(*object.Str).Value

	return &object.Str{Value: strings.TrimRight(val, chars)}, nil
}

// strTrimLeftFunc returns a string with trailing whitespace removed from the left
func strTrimLeftFunc(receiver object.Object, args ...object.Object) (object.Object, error) {
	chars := defaultCharTrim
	if len(args) > 0 {
		str, ok := args[0].(*object.Str)

		if !ok {
			msg := fmt.Sprintf(fail.ErrFuncFirstArgStr, object.STR_OBJ, "trimLeft")
			return nil, errors.New(msg)
		}

		chars = str.Value
	}

	val := receiver.(*object.Str).Value

	return &object.Str{Value: strings.TrimLeft(val, chars)}, nil
}

// strRepeatFunc returns a string repeated n times
func strRepeatFunc(receiver object.Object, args ...object.Object) (object.Object, error) {
	if len(args) == 0 {
		msg := fmt.Sprintf(fail.ErrFuncMissingArg, object.STR_OBJ, "repeat")
		return nil, errors.New(msg)
	}

	firstArg, ok := args[0].(*object.Int)
	if !ok {
		msg := fmt.Sprintf(fail.ErrFuncFirstArgInt, object.STR_OBJ, "repeat")
		return nil, errors.New(msg)
	}

	val := receiver.(*object.Str).Value
	repeated := strings.Repeat(val, int(firstArg.Value))

	return &object.Str{Value: repeated}, nil
}
