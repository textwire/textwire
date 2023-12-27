package object

import (
	"errors"
)

type Env struct {
	store map[string]Object
}

func NewEnv() *Env {
	s := make(map[string]Object)
	return &Env{store: s}
}

func EnvFromMap(m map[string]interface{}) (*Env, error) {
	env := NewEnv()

	for key, val := range m {
		switch v := val.(type) {
		case string:
			env.Set(key, &String{Value: v})
		case bool:
			env.Set(key, &Boolean{Value: v})
		case float32:
			env.Set(key, &Float32{Value: v})
		case float64:
			env.Set(key, &Float64{Value: v})
		case int:
			env.Set(key, &Int{Value: v})
		case int8:
			env.Set(key, &Int8{Value: v})
		case int16:
			env.Set(key, &Int16{Value: v})
		case int32:
			env.Set(key, &Int32{Value: v})
		case int64:
			env.Set(key, &Int64{Value: v})
		case uint:
			env.Set(key, &Uint{Value: v})
		case uint8:
			env.Set(key, &Uint8{Value: v})
		case uint16:
			env.Set(key, &Uint16{Value: v})
		case uint32:
			env.Set(key, &Uint32{Value: v})
		case uint64:
			env.Set(key, &Uint64{Value: v})
		default:
			return nil, errors.New("Unsupported type for Textwire parser")
		}
	}

	return env, nil
}

func (e *Env) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	return obj, ok
}

func (e *Env) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
