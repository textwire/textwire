package object

import (
	"errors"
	"fmt"

	fail "github.com/textwire/textwire/v2/fail"
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
		obj := NativeToObject(val)

		if obj == nil {
			return nil, fail.New(0, "", "template",
				fail.ErrUnsupportedType, val)
		}

		err := env.Set(key, obj)
		if err != nil {
			return nil, fail.New(0, "", "evaluator", "%s", err.Error())
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

func (e *Env) Set(key string, val Object) error {
	if key == "loop" {
		return errors.New(fail.ErrLoopVariableIsReserved)
	}

	if oldVar, ok := e.isTypeMismatch(key, val); ok {
		return e.variableMismatchError(key, oldVar, val)
	}

	e.store[key] = val

	return nil
}

func (e *Env) SetLoopVar(pairs map[string]Object) {
	e.store["loop"] = &Obj{Pairs: pairs}
}

func (e *Env) isTypeMismatch(key string, val Object) (Object, bool) {
	oldVar, ok := e.Get(key)
	return oldVar, (ok && oldVar != nil && oldVar.Type() != val.Type())
}

func (e *Env) variableMismatchError(key string, oldVar, val Object) error {
	msg := fmt.Sprintf(fail.ErrVariableTypeMismatch, key, oldVar.Type(), val.Type())
	return errors.New(msg)
}
