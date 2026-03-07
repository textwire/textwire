package config

type StringCustomFunc func(s string, args ...any) any
type ArrayCustomFunc func(a []any, args ...any) any
type IntegerCustomFunc func(i int, args ...any) any
type FloatCustomFunc func(f float64, args ...any) any
type BooleanCustomFunc func(b bool, args ...any) any
type MapCustomFunc func(o map[string]any, args ...any) any

type Func struct {
	String  map[string]StringCustomFunc
	Array   map[string]ArrayCustomFunc
	Integer map[string]IntegerCustomFunc
	Float   map[string]FloatCustomFunc
	Boolean map[string]BooleanCustomFunc
	Map     map[string]MapCustomFunc
}

func NewFunc() *Func {
	return &Func{
		String:  map[string]StringCustomFunc{},
		Array:   map[string]ArrayCustomFunc{},
		Integer: map[string]IntegerCustomFunc{},
		Float:   map[string]FloatCustomFunc{},
		Boolean: map[string]BooleanCustomFunc{},
		Map:     map[string]MapCustomFunc{},
	}
}
