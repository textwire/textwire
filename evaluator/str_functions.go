package evaluator

import (
	"github.com/textwire/textwire/object"
)

// strLenFunc returns the length of the given string
func strLenFunc(receiver object.Object, args ...object.Object) object.Object {
	str := receiver.(*object.Str)
	val := len(str.Value)
	return &object.Int{Value: int64(val)}
}
