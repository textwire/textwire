package evaluator

import "github.com/textwire/textwire/v2/config"

type EvalContext struct {
	absPath string
	conf    *config.Config
}

func NewContext(absPath string, conf *config.Config) *EvalContext {
	return &EvalContext{
		absPath: absPath,
		conf:    conf,
	}
}
