package object

import "errors"

type Env struct {
	store map[string]Object
}

func NewEnv() *Env {
	s := make(map[string]Object)
	return &Env{store: s}
}

func EnvFromMap(m map[string]interface{}) (*Env, error) {
	var env *Env

	for key, val := range m {
		switch val.(type) {
		case string:
			env.Set(key, &String{Value: val.(string)})
		case int:
			env.Set(key, &Integer{Value: int64(val.(int))})
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
