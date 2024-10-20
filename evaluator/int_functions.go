package evaluator

import "github.com/textwire/textwire/v2/object"

func intFloatFunc(receiver object.Object, args ...object.Object) object.Object {
	floatVal := receiver.(*object.Int).Value
	return &object.Float{Value: float64(floatVal)}
}
