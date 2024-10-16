package option

import "github.com/textwire/textwire/v2/object"

type StrCustomFunc func(s *object.Str, args ...object.Object) object.Object
type ArrayCustomFunc func(a *object.Array, args ...object.Object) object.Object
type IntCustomFunc func(i *object.Int, args ...object.Object) object.Object
type FloatCustomFunc func(f *object.Float, args ...object.Object) object.Object
type BoolCustomFunc func(b *object.Bool, args ...object.Object) object.Object

type Func struct {
	Str   map[string]StrCustomFunc
	Arr   map[string]ArrayCustomFunc
	Int   map[string]IntCustomFunc
	Float map[string]FloatCustomFunc
	Bool  map[string]BoolCustomFunc
}

func NewFunc() *Func {
	return &Func{
		Str:   make(map[string]StrCustomFunc),
		Arr:   make(map[string]ArrayCustomFunc),
		Int:   make(map[string]IntCustomFunc),
		Float: make(map[string]FloatCustomFunc),
		Bool:  make(map[string]BoolCustomFunc),
	}
}
