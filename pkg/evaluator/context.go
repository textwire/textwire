package evaluator

import "github.com/textwire/textwire/v4/pkg/value"

// SlotsStore stores slots. Example below:
//
//	SlotsStore{
//	    "components/book": {
//		    "title": value.Value{},
//		    "author": value.Value{},
//	    },
//	}
type SlotsStore = map[string]map[string]value.Value

// Context is evaluator context that is being passed through all
// the evaluator values to carry scope and path to the current file.
type Context struct {
	scope   *value.Scope // current value's scope
	absPath string       // absolute path to the file being executed

	// inserts should be used inside of layouts.
	// - key is the name of the insert.
	// - value is evaluated ASTs into value.
	inserts map[string]value.Value

	// slots should be used inside component files.
	slots SlotsStore

	// compSlots is used during component block evaluation to redirect
	// @provide and @provideif output to the component's context instead
	// of the current context. This allows nested provides within component
	// blocks to populate the component's slots.
	compSlots SlotsStore
}

func NewContext(scope *value.Scope, absPath string) *Context {
	if scope == nil {
		panic("scope should never be nil in evaluator context")
	}

	return &Context{
		scope:   scope,
		absPath: absPath,
		inserts: map[string]value.Value{},
		slots:   SlotsStore{},
	}
}
