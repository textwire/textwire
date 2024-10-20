package evaluator

import (
	"github.com/textwire/textwire/v2/object"
)

// arrayLenFunc returns the length of the given array
func arrayLenFunc(receiver object.Object, args ...object.Object) (object.Object, error) {
	length := len(receiver.(*object.Array).Elements)
	return &object.Int{Value: int64(length)}, nil
}

// arrayJoinFunc joins the elements of the given array with the given separator
func arrayJoinFunc(receiver object.Object, args ...object.Object) (object.Object, error) {
	var separator string

	if len(args) == 0 {
		separator = ","
	} else {
		separator = args[0].(*object.Str).Value
	}

	elements := receiver.(*object.Array).Elements

	var result string

	for i, el := range elements {
		if i > 0 {
			result += separator
		}
		result += el.String()
	}

	return &object.Str{Value: result}, nil
}

// arrayRandFunc returns a random element from the given array
func arrayRandFunc(receiver object.Object, args ...object.Object) (object.Object, error) {
	elements := receiver.(*object.Array).Elements
	length := len(elements)

	if length == 0 {
		return &object.Nil{}, nil
	}

	return elements[0], nil
}

// arrayReverseFunc reverses the elements of the given array
func arrayReverseFunc(receiver object.Object, args ...object.Object) (object.Object, error) {
	elements := receiver.(*object.Array).Elements
	length := len(elements)

	if length == 0 {
		return receiver, nil
	}

	reversed := make([]object.Object, length)

	for i, el := range elements {
		reversed[length-i-1] = el
	}

	return &object.Array{Elements: reversed}, nil
}
