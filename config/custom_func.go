package config

type StrCustomFunc func(s string, args ...any) any
type ArrayCustomFunc func(a []any, args ...any) any
type IntCustomFunc func(i int, args ...any) any
type FloatCustomFunc func(f float64, args ...any) any
type BoolCustomFunc func(b bool, args ...any) any
type ObjectCustomFunc func(o map[string]any, args ...any) any

type Func struct {
	Str   map[string]StrCustomFunc
	Arr   map[string]ArrayCustomFunc
	Int   map[string]IntCustomFunc
	Float map[string]FloatCustomFunc
	Bool  map[string]BoolCustomFunc
	Obj   map[string]ObjectCustomFunc
}

func NewFunc() *Func {
	return &Func{
		Str:   map[string]StrCustomFunc{},
		Arr:   map[string]ArrayCustomFunc{},
		Int:   map[string]IntCustomFunc{},
		Float: map[string]FloatCustomFunc{},
		Bool:  map[string]BoolCustomFunc{},
		Obj:   map[string]ObjectCustomFunc{},
	}
}
