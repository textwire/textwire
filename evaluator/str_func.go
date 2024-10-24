package evaluator

import (
	"errors"
	"fmt"
	"html"
	"strings"

	"github.com/textwire/textwire/v2/fail"
	"github.com/textwire/textwire/v2/object"
)

// strLenFunc returns the length of the given string
func strLenFunc(receiver object.Object, _ ...object.Object) (object.Object, error) {
	val := receiver.(*object.Str).Value
	return &object.Int{Value: int64(len(val))}, nil
}

// strSplitFunc returns a list of strings split by the given separator
func strSplitFunc(receiver object.Object, args ...object.Object) (object.Object, error) {
	separator := " "

	if len(args) > 0 {
		str, ok := args[0].(*object.Str)

		if !ok {
			msg := fmt.Sprintf(fail.ErrFuncFirstArgStr, "split", "string")
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
	chars := "\t \n\r"

	if len(args) > 0 {
		str, ok := args[0].(*object.Str)

		if !ok {
			msg := fmt.Sprintf(fail.ErrFuncFirstArgStr, "trim", "string")
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
		msg := fmt.Sprintf(fail.ErrFuncRequiresOneArg, "contains", "string")
		return nil, errors.New(msg)
	}

	substr, ok := args[0].(*object.Str)

	if !ok {
		msg := fmt.Sprintf(fail.ErrFuncFirstArgStr, "contains", "string")
		return nil, errors.New(msg)
	}

	val := receiver.(*object.Str).Value

	return &object.Bool{Value: strings.Contains(val, substr.Value)}, nil
}
