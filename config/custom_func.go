package config

type StrCustomFunc func(s string, args ...any) string
type ArrayCustomFunc func(a []any, args ...any) []any
type IntCustomFunc func(i int, args ...any) int
type FloatCustomFunc func(f float64, args ...any) float64
type BoolCustomFunc func(b bool, args ...any) bool

type Func struct {
	Str   map[string]StrCustomFunc
	Arr   map[string]ArrayCustomFunc
	Int   map[string]IntCustomFunc
	Float map[string]FloatCustomFunc
	Bool  map[string]BoolCustomFunc
}

func NewFunc() *Func {
	return &Func{
		Str:   map[string]StrCustomFunc{},
		Arr:   map[string]ArrayCustomFunc{},
		Int:   map[string]IntCustomFunc{},
		Float: map[string]FloatCustomFunc{},
		Bool:  map[string]BoolCustomFunc{},
	}
}
