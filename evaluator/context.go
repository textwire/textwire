package evaluator

import "github.com/textwire/textwire/v2/option"

type EvalContext struct {
	absPath    string
	customFunc *option.Func
}

func NewContext(absPath string, customFunc *option.Func) *EvalContext {
	return &EvalContext{
		absPath:    absPath,
		customFunc: customFunc,
	}
}
