package evaluator

import "github.com/textwire/textwire/object"

// arrayLenFunc returns the length of the given array
func arrayLenFunc(receiver object.Object, args ...object.Object) object.Object {
	length := len(receiver.(*object.Array).Elements)
	return &object.Int{Value: int64(length)}
}
