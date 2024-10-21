package config

type StrCustomFunc func(s string, args ...interface{}) string
type ArrayCustomFunc func(a []interface{}, args ...interface{}) []interface{}
type IntCustomFunc func(i int, args ...interface{}) int
type FloatCustomFunc func(f float64, args ...interface{}) float64
type BoolCustomFunc func(b bool, args ...interface{}) bool

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