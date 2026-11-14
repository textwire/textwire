package value

import (
	"errors"
	"fmt"

	"codeberg.org/textwire/textwire/v6/pkg/fail"
)

type Scope struct {
	vars   map[string]Literal
	parent *Scope
}

// NewScope returns a scope with global value.
// A nil global falls back to the shared empty one.
func NewScope(global *Obj) *Scope {
	if global == nil {
		global = emptyGlobal
	}

	return &Scope{vars: map[string]Literal{"global": global}}
}

func NewChildScope(parent *Scope) *Scope {
	return &Scope{vars: map[string]Literal{}, parent: parent}
}

func NewScopeFromMap(data map[string]any, global *Obj) (*Scope, *fail.Error) {
	scope := NewScope(global)

	for key, val := range data {
		obj := NativeToValue(val)
		if errObj, ok := obj.(*Error); ok {
			return nil, errObj.Err
		}

		if err := scope.Set(key, obj); err != nil {
			return nil, fail.New(nil, "", fail.OriginEval, "%s", err.Error())
		}
	}

	return scope, nil
}

func (e *Scope) Get(name string) (Literal, bool) {
	if val, ok := e.vars[name]; ok {
		return val, true
	}

	if e.parent != nil {
		return e.parent.Get(name)
	}

	return nil, false
}

func (e *Scope) Set(key string, val Literal) error {
	if key == "loop" || key == "global" {
		return errors.New(fail.ErrReservedIdentifiers)
	}

	if oldVar, ok := e.isTypeMismatch(key, val); ok {
		return e.identifierMismatchError(key, oldVar, val)
	}

	writeScope := e.scopeOf(key)
	if writeScope == nil {
		writeScope = e
	}

	writeScope.vars[key] = val

	return nil
}

// scopeOf walks up from the current scope and returns the closest one
// that defines key. It returns nil when no scope in the chain defines it.
func (e *Scope) scopeOf(key string) *Scope {
	if _, ok := e.vars[key]; ok {
		return e
	}

	if e.parent == nil {
		return nil
	}

	return e.parent.scopeOf(key)
}

func (e *Scope) SetLoopVariable(key string, val Literal) error {
	if key == "loop" || key == "global" {
		return errors.New(fail.ErrReservedIdentifiers)
	}

	e.vars[key] = val
	return nil
}

func (e *Scope) SetLoopVar(pairs map[string]Literal) {
	e.vars["loop"] = NewObj(pairs)
}

func (e *Scope) isTypeMismatch(key string, val Value) (Value, bool) {
	oldVar, ok := e.Get(key)
	return oldVar, (ok && oldVar != nil && oldVar.Type() != val.Type())
}

func (e *Scope) identifierMismatchError(key string, oldVar, val Value) error {
	msg := fmt.Sprintf(fail.ErrIdentTypeMismatch, key, oldVar.Type(), val.Type())
	return errors.New(msg)
}
