package evaluator

import "github.com/textwire/textwire/v2/config"

type EvalContext struct {
	absPath    string
	customFunc *config.Func
}

func NewContext(absPath string, customFunc *config.Func) *EvalContext {
	return &EvalContext{
		absPath:    absPath,
		customFunc: customFunc,
	}
}
