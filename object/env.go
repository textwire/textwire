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
		obj := goValueToObj(val)

		if obj == nil {
			return nil, errors.New("Unsupported type for Textwire language")
		}

		env.Set(key, obj)
	}

	return env, nil
}

func goValueToObj(val interface{}) Object {
	switch v := val.(type) {
	case string:
		return &String{Value: v}
	case bool:
		return &Boolean{Value: v}
	case float32:
		return &Float{Value: float64(v)}
	case float64:
		return &Float{Value: v}
	case int64:
		return &Int{Value: v}
	case int:
		return &Int{Value: int64(v)}
	case int8:
		return &Int{Value: int64(v)}
	case int16:
		return &Int{Value: int64(v)}
	case int32:
		return &Int{Value: int64(v)}
	case uint:
		return &Int{Value: int64(v)}
	case uint8:
		return &Int{Value: int64(v)}
	case uint16:
		return &Int{Value: int64(v)}
	case uint32:
		return &Int{Value: int64(v)}
	case uint64:
		return &Int{Value: int64(v)}
	default:
		return nil
	}
}

func (e *Env) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	return obj, ok
}

func (e *Env) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
