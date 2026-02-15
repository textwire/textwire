package evaluator

import "github.com/textwire/textwire/v3/pkg/object"

// SlotsStore stores slots. Example below:
//
//	SlotsStore{
//	    "components/book": {
//	        "title": object.Object{},
//	        "author": object.Object{},
//	    },
//	}
type SlotsStore = map[string]map[string]object.Object

// Context is evaluator context that is being passed through all
// the evaluator objects to carry scope and path to the current file.
type Context struct {
	scope   *object.Scope // current object's scope
	absPath string        // absolute path to the file being executed

	// inserts should be used inside of layouts.
	// - key is the name of the insert.
	// - value is evaluated ASTs into object.
	inserts map[string]object.Object

	// slots should be used inside component files.
	slots SlotsStore
}

func NewContext(scope *object.Scope, absPath string) *Context {
	if scope == nil {
		panic("scope should never be nil in evaluator context")
	}

	return &Context{
		scope:   scope,
		absPath: absPath,
		inserts: map[string]object.Object{},
		slots:   SlotsStore{},
	}
}
