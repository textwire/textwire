package evaluator

import "github.com/textwire/textwire/object"

// arrayLenFunc returns the length of the given array
func arrayLenFunc(receiver object.Object, args ...object.Object) object.Object {
	length := len(receiver.(*object.Array).Elements)
	return &object.Int{Value: int64(length)}
}

// arrayJoinFunc joins the elements of the given array with the given separator
func arrayJoinFunc(receiver object.Object, args ...object.Object) object.Object {
	separator := args[0].(*object.Str).Value
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
