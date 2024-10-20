package evaluator

import "github.com/textwire/textwire/v2/object"

// arrayLenFunc returns the length of the given array
func arrayLenFunc(receiver object.Object, args ...object.Object) object.Object {
	length := len(receiver.(*object.Array).Elements)
	return &object.Int{Value: int64(length)}
}

// arrayJoinFunc joins the elements of the given array with the given separator
func arrayJoinFunc(receiver object.Object, args ...object.Object) object.Object {
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

	return &object.Str{Value: result}
}

func arrayRandFunc(receiver object.Object, args ...object.Object) object.Object {
	elements := receiver.(*object.Array).Elements
	length := len(elements)

	if length == 0 {
		return &object.Nil{}
	}

	return elements[0]
}

func arrayReverseFunc(receiver object.Object, args ...object.Object) object.Object {
	elements := receiver.(*object.Array).Elements
	length := len(elements)

	if length == 0 {
		return receiver
	}

	reversed := make([]object.Object, length)

	for i, el := range elements {
		reversed[length-i-1] = el
	}

	return &object.Array{Elements: reversed}
}
