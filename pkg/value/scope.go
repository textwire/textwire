package value

import (
	"errors"
	"fmt"

	"github.com/textwire/textwire/v3/pkg/fail"
)

type Scope struct {
	vars   map[string]Literal
	parent *Scope
}

func NewScope() *Scope {
	vars := map[string]Literal{}
	vars["global"] = NewObj(nil)

	return &Scope{vars: vars}
}

func NewScopeFromMap(data map[string]any) (*Scope, *fail.Error) {
	scope := NewScope()

	for key, val := range data {
		obj := NativeToValue(val)
		if obj == nil {
			return nil, fail.New(nil, "", "template", fail.ErrUnsupportedType, val)
		}

		if err := scope.Set(key, obj); err != nil {
			return nil, fail.New(nil, "", "evaluator", "%s", err.Error())
		}
	}

	return scope, nil
}

func (s *Scope) Child() *Scope {
	child := NewScope()
	child.parent = s
	return child
}

func (e *Scope) Get(name string) (Literal, bool) {
	obj, ok := e.vars[name]
	if !ok && e.parent != nil {
		obj, ok = e.parent.Get(name)
	}
	return obj, ok
}

func (e *Scope) Set(key string, val Literal) error {
	if key == "loop" || key == "global" {
		return errors.New(fail.ErrReservedIdentifiers)
	}

	if oldVar, ok := e.isTypeMismatch(key, val); ok {
		return e.identifierMismatchError(key, oldVar, val)
	}

	e.vars[key] = val

	return nil
}

func (e *Scope) SetLoopVar(pairs map[string]Literal) {
	e.vars["loop"] = NewObj(pairs)
}

func (e *Scope) AddGlobal(key string, val any) {
	var globalObj *Obj

	switch v := e.vars["global"].(type) {
	case *Obj:
		globalObj = v
	case nil:
		globalObj = NewObj(nil)
		e.vars["global"] = globalObj
	default:
		globalObj = NewObj(nil)
		e.vars["global"] = globalObj
	}

	// Ensure Pairs map is initialized
	if globalObj.Pairs == nil {
		globalObj.Pairs = map[string]Literal{}
	}

	globalObj.Pairs[key] = NativeToValue(val)
}

func (e *Scope) isTypeMismatch(key string, val Value) (Value, bool) {
	oldVar, ok := e.Get(key)
	return oldVar, (ok && oldVar != nil && oldVar.Type() != val.Type())
}

func (e *Scope) identifierMismatchError(key string, oldVar, val Value) error {
	msg := fmt.Sprintf(fail.ErrIdentifierTypeMismatch, key, oldVar.Type(), val.Type())
	return errors.New(msg)
}
