package evaluator

import (
	"github.com/textwire/textwire/v3/pkg/object"
)

// objCamelFunc converts object keys to camel case recursively
func objCamelFunc(receiver object.Object, _ ...object.Object) (object.Object, error) {
	obj := receiver.(*object.Obj)
	return &object.Obj{Pairs: obj.ToCamel()}, nil
}
