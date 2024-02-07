package evaluator

import (
	"strings"

	"github.com/textwire/textwire/object"
)

// strLenFunc returns the length of the given string
func strLenFunc(receiver object.Object, args ...object.Object) object.Object {
	str := receiver.(*object.Str)
	val := len(str.Value)
	return &object.Int{Value: int64(val)}
}

// strSplitFunc returns a list of strings split by the given separator
func strSplitFunc(receiver object.Object, args ...object.Object) object.Object {
	separator := " "

	if len(args) > 0 && args[0].Type() == object.STR_OBJ {
		separator = args[0].(*object.Str).Value
	}

	str := receiver.(*object.Str)
	stringItems := strings.Split(str.Value, separator)

	var elems []object.Object

	for _, val := range stringItems {
		elems = append(elems, &object.Str{Value: val})
	}

	return &object.Array{Elements: elems}
}
