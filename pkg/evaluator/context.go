package evaluator

import (
	"github.com/textwire/textwire/v6/pkg/ast"
	"github.com/textwire/textwire/v6/pkg/value"
)

// SlotsStore stores slots. Example below:
//
//	SlotsStore{
//	    "components/book": {
//		    "title": value.Value{},
//		    "author": value.Value{},
//	    },
//	}
type SlotsStore = map[string]map[string]value.Value

// useState is shared by a template's child contexts, so @use directive that
// runs inside a block is visible when the template finishes evaluating.
type useState struct {
	dir *ast.UseDir
	ran bool
}

// Context is evaluator context that is being passed through all
// the evaluator values to carry scope and path to the current file.
type Context struct {
	scope   *value.Scope // current value's scope
	absPath string       // absolute path to the file being executed

	// inserts should be used inside of layouts.
	// - key is the name of the insert.
	// - value is evaluated ASTs into value.
	inserts map[string]value.Value

	useState *useState

	// slots should be used inside component files.
	slots SlotsStore
}

func NewContext(scope *value.Scope, absPath string) *Context {
	if scope == nil {
		panic("scope should never be nil in evaluator context")
	}

	return &Context{
		scope:    scope,
		absPath:  absPath,
		inserts:  map[string]value.Value{},
		useState: &useState{},
		slots:    SlotsStore{},
	}
}

func (ctx *Context) childContext() *Context {
	return &Context{
		scope:    ctx.scope,
		absPath:  ctx.absPath,
		inserts:  ctx.inserts,
		useState: ctx.useState,
		slots:    ctx.slots,
	}
}

func (ctx *Context) childContextWithNewScope() *Context {
	childCtx := ctx.childContext()
	childCtx.scope = value.NewChildScope(ctx.scope)
	return childCtx
}
