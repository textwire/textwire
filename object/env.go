package object

import (
	"errors"
	"fmt"

	fail "github.com/textwire/textwire/v3/fail"
)

type Env struct {
	store map[string]Object
	outer *Env
}

func NewEnv() *Env {
	store := map[string]Object{}
	store["global"] = NewObj(nil)

	return &Env{store: store}
}

func NewEnclosedEnv(outer *Env) *Env {
	env := NewEnv()
	env.outer = outer
	return env
}

func EnvFromMap(data map[string]any) (*Env, *fail.Error) {
	env := NewEnv()

	for key, val := range data {
		obj := NativeToObject(val)

		if obj == nil {
			return nil, fail.New(0, "", "template", fail.ErrUnsupportedType, val)
		}

		if err := env.Set(key, obj); err != nil {
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
	if key == "loop" || key == "global" {
		return errors.New(fail.ErrReservedIdentifiers)
	}

	if oldVar, ok := e.isTypeMismatch(key, val); ok {
		return e.identifierMismatchError(key, oldVar, val)
	}

	e.store[key] = val

	return nil
}

func (e *Env) SetLoopVar(pairs map[string]Object) {
	e.store["loop"] = NewObj(pairs)
}

func (e *Env) AddGlobalData(key string, val any) {
	var globalObj *Obj

	switch v := e.store["global"].(type) {
	case *Obj:
		globalObj = v
	case nil:
		globalObj = NewObj(nil)
		e.store["global"] = globalObj
	default:
		globalObj = NewObj(nil)
		e.store["global"] = globalObj
	}

	// Ensure Pairs map is initialized
	if globalObj.Pairs == nil {
		globalObj.Pairs = map[string]Object{}
	}

	globalObj.Pairs[key] = NativeToObject(val)
}

func (e *Env) isTypeMismatch(key string, val Object) (Object, bool) {
	oldVar, ok := e.Get(key)
	return oldVar, (ok && oldVar != nil && oldVar.Type() != val.Type())
}

func (e *Env) identifierMismatchError(key string, oldVar, val Object) error {
	msg := fmt.Sprintf(fail.ErrIdentifierTypeMismatch, key, oldVar.Type(), val.Type())
	return errors.New(msg)
}
