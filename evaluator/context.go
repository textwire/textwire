package evaluator

import "github.com/textwire/textwire/v2/config"

type EvalContext struct {
	absPath string
	// TODO: use this field
	config *config.Config
}

func NewContext(absPath string) *EvalContext {
	return &EvalContext{absPath: absPath}
}
