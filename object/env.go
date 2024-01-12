package object

import (
	"errors"
)

type Env struct {
	store map[string]Object
	outer *Env
}

func NewEnv() *Env {
	s := make(map[string]Object)
	return &Env{store: s}
}

func NewEnclosedEnv(outer *Env) *Env {
	env := NewEnv()
	env.outer = outer
	return env
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
			env.Set(key, &Float{Value: float64(v)})
		case float64:
			env.Set(key, &Float{Value: v})
		case int64:
			env.Set(key, &Int{Value: v})
		case int:
			env.Set(key, &Int{Value: int64(v)})
		case int8:
			env.Set(key, &Int{Value: int64(v)})
		case int16:
			env.Set(key, &Int{Value: int64(v)})
		case int32:
			env.Set(key, &Int{Value: int64(v)})
		case uint:
			env.Set(key, &Int{Value: int64(v)})
		case uint8:
			env.Set(key, &Int{Value: int64(v)})
		case uint16:
			env.Set(key, &Int{Value: int64(v)})
		case uint32:
			env.Set(key, &Int{Value: int64(v)})
		case uint64:
			env.Set(key, &Int{Value: int64(v)})
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
