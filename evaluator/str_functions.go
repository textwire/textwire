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
	str := receiver.(*object.Str).Value
	return &object.Int{Value: int64(len(str))}, nil
}

// strSplitFunc returns a list of strings split by the given separator
func strSplitFunc(receiver object.Object, args ...object.Object) (object.Object, error) {
	separator := " "

	if len(args) > 0 {
		str, ok := args[0].(*object.Str)

		if !ok {
			msg := fmt.Sprintf(fail.ErrFuncFirstArgStr, "split")
			return nil, errors.New(msg)
		}

		separator = str.Value
	}

	str := receiver.(*object.Str).Value
	stringItems := strings.Split(str, separator)

	var elems []object.Object

	for _, val := range stringItems {
		elems = append(elems, &object.Str{Value: val})
	}

	return &object.Array{Elements: elems}, nil
}

// strRawFunc prevents escaping HTML tags in a string
func strRawFunc(receiver object.Object, _ ...object.Object) (object.Object, error) {
	str := receiver.(*object.Str)
	return &object.Str{Value: html.UnescapeString(str.Value)}, nil
}

// strTrimFunc returns a string with leading and trailing whitespace removed
func strTrimFunc(receiver object.Object, args ...object.Object) (object.Object, error) {
	chars := "\t \n\r"

	if len(args) > 0 {
		str, ok := args[0].(*object.Str)

		if !ok {
			msg := fmt.Sprintf(fail.ErrFuncFirstArgStr, "trim")
			return nil, errors.New(msg)
		}

		chars = str.Value
	}

	str := receiver.(*object.Str).Value

	return &object.Str{Value: strings.Trim(str, chars)}, nil
}

// strUpperFunc returns a string with all characters in uppercase
func strUpperFunc(receiver object.Object, _ ...object.Object) (object.Object, error) {
	str := receiver.(*object.Str)
	return &object.Str{Value: strings.ToUpper(str.Value)}, nil
}

// strLowerFunc returns a string with all characters in lowercase
func strLowerFunc(receiver object.Object, _ ...object.Object) (object.Object, error) {
	str := receiver.(*object.Str)
	return &object.Str{Value: strings.ToLower(str.Value)}, nil
}
