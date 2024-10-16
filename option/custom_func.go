package option

import "github.com/textwire/textwire/v2/object"

type Func struct {
	Str   map[string]object.BuiltinFunction
	Arr   map[string]object.BuiltinFunction
	Int   map[string]object.BuiltinFunction
	Float map[string]object.BuiltinFunction
	Bool  map[string]object.BuiltinFunction
}

func NewFunc() *Func {
	return &Func{
		Str:   make(map[string]object.BuiltinFunction),
		Arr:   make(map[string]object.BuiltinFunction),
		Int:   make(map[string]object.BuiltinFunction),
		Float: make(map[string]object.BuiltinFunction),
		Bool:  make(map[string]object.BuiltinFunction),
	}
}
