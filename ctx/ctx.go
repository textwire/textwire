package ctx

import "github.com/textwire/textwire/v2/config"

type EvalCtx struct {
	AbsPath    string
	CustomFunc *config.Func
	Config     *config.Config
}

func NewContext(absPath string, customFunc *config.Func, conf *config.Config) *EvalCtx {
	return &EvalCtx{
		AbsPath:    absPath,
		CustomFunc: customFunc,
		Config:     conf,
	}
}
