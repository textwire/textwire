package evaluator

import "github.com/textwire/textwire/v3/object"

// Context is evaluator context that is being passed through all
// the evaluator objects to carry scope and path to the current file.
type Context struct {
	scope   *object.Scope // current object's scope
	absPath string        // absolute path to the file being executed
}

func NewContext(scope *object.Scope, absPath string) *Context {
	if scope == nil {
		panic("scope should never be nil in evaluator context")
	}

	return &Context{scope: scope, absPath: absPath}
}
