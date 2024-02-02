package evaluator

type EvalContext struct {
	absPath string
}

func NewContext(absPath string) *EvalContext {
	return &EvalContext{absPath: absPath}
}
