package evaluator

import (
	"html"
	"strings"

	"github.com/textwire/textwire/v2/object"
)

// strLenFunc returns the length of the given string
func strLenFunc(receiver object.Object, _ ...object.Object) object.Object {
	str := receiver.(*object.Str).Value
	return &object.Int{Value: int64(len(str))}
}

// strSplitFunc returns a list of strings split by the given separator
func strSplitFunc(receiver object.Object, args ...object.Object) object.Object {
	separator := " "

	if len(args) > 0 && args[0].Type() == object.STR_OBJ {
		separator = args[0].(*object.Str).Value
	}

	str := receiver.(*object.Str).Value
	stringItems := strings.Split(str, separator)

	var elems []object.Object

	for _, val := range stringItems {
		elems = append(elems, &object.Str{Value: val})
	}

	return &object.Array{Elements: elems}
}

// strRawFunc prevents escaping HTML tags in a string
func strRawFunc(receiver object.Object, _ ...object.Object) object.Object {
	str := receiver.(*object.Str)
	return &object.Str{Value: html.UnescapeString(str.Value)}
}

// strTrimFunc returns a string with leading and trailing whitespace removed
func strTrimFunc(receiver object.Object, args ...object.Object) object.Object {
	chars := "\t \n\r"

	if len(args) > 0 && args[0].Type() == object.STR_OBJ {
		chars = args[0].(*object.Str).Value
	}

	str := receiver.(*object.Str).Value

	return &object.Str{Value: strings.Trim(str, chars)}
}

// strUpperFunc returns a string with all characters in uppercase
func strUpperFunc(receiver object.Object, _ ...object.Object) object.Object {
	str := receiver.(*object.Str)
	return &object.Str{Value: strings.ToUpper(str.Value)}
}

// strLowerFunc returns a string with all characters in lowercase
func strLowerFunc(receiver object.Object, _ ...object.Object) object.Object {
	str := receiver.(*object.Str)
	return &object.Str{Value: strings.ToLower(str.Value)}
}
