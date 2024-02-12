package object

import (
	"errors"

	"github.com/textwire/textwire/fail"
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

func EnvFromMap(data map[string]interface{}) (*Env, *fail.Error) {
	env := NewEnv()

	for key, val := range data {
		obj := nativeToObject(val)

		if obj == nil {
			return nil, fail.New(0, "", "template", fail.ErrUnsupportedType, val)
		}

		err := env.Set(key, obj)

		if err != nil {
			return nil, fail.New(0, "", "template", err.Error())
		}
	}

	return env, nil
}

func (e *Env) Get(name string) (Object, bool) {
	obj, ok := e.store[name]

	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}

	return obj, ok
}

func (e *Env) Set(name string, val Object) error {
	if name == "loop" {
		return errors.New(fail.ErrLoopVariableIsReserved)
	}

	e.store[name] = val

	return nil
}
