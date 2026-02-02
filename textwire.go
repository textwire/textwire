package textwire

import (
	"github.com/textwire/textwire/v3/config"
	"github.com/textwire/textwire/v3/evaluator"
	"github.com/textwire/textwire/v3/fail"
	"github.com/textwire/textwire/v3/object"
)

var (
	userConfig = config.New("templates", ".tw", "", false)
	customFunc = config.NewFunc()
)

func NewTemplate(opt *config.Config) (*Template, error) {
	Configure(opt)

	twFiles, err := findTwFiles()
	if err != nil {
		return nil, fail.FromError(err, 0, "", "template").Error()
	}

	parseErr := parsePrograms(twFiles)
	if parseErr != nil {
		return nil, parseErr.Error()
	}

	return &Template{twFiles: twFiles}, nil
}

func EvaluateString(inp string, data map[string]any) (string, error) {
	prog, errs := parseStr(inp)
	if len(errs) != 0 {
		return "", errs[0].Error()
	}

	env, err := object.EnvFromMap(data)
	if err != nil {
		return "", err.Error()
	}

	eval := evaluator.New(customFunc, nil)

	evaluated := eval.Eval(prog, env, prog.Filepath)
	if evaluated.Is(object.ERR_OBJ) {
		return "", evaluated.(*object.Error).Err.Error()
	}

	return evaluated.String(), nil
}

func EvaluateFile(absPath string, data map[string]any) (string, error) {
	twFile := NewTwFile("", absPath)

	content, err := fileContent(twFile)
	if err != nil {
		return "", fail.FromError(err, 0, absPath, "template").Error()
	}

	result, err := EvaluateString(content, data)
	if err != nil {
		return "", err
	}

	return result, nil
}

func RegisterStrFunc(name string, fn config.StrCustomFunc) error {
	if _, ok := customFunc.Str[name]; ok {
		return fail.New(0, "", "API", fail.ErrFuncAlreadyDefined,
			name, "strings").Error()
	}

	customFunc.Str[name] = fn

	return nil
}

func RegisterArrFunc(name string, fn config.ArrayCustomFunc) error {
	if _, ok := customFunc.Arr[name]; ok {
		return fail.New(0, "", "API", fail.ErrFuncAlreadyDefined,
			name, "arrays").Error()
	}

	customFunc.Arr[name] = fn

	return nil
}

func RegisterObjFunc(name string, fn config.ObjectCustomFunc) error {
	if _, ok := customFunc.Obj[name]; ok {
		return fail.New(0, "", "API", fail.ErrFuncAlreadyDefined,
			name, "objects").Error()
	}

	customFunc.Obj[name] = fn

	return nil
}

func RegisterIntFunc(name string, fn config.IntCustomFunc) error {
	if _, ok := customFunc.Int[name]; ok {
		return fail.New(0, "", "API", fail.ErrFuncAlreadyDefined,
			name, "integers").Error()
	}

	customFunc.Int[name] = fn

	return nil
}

func RegisterFloatFunc(name string, fn config.FloatCustomFunc) error {
	if _, ok := customFunc.Float[name]; ok {
		return fail.New(0, "", "API", fail.ErrFuncAlreadyDefined,
			name, "floats").Error()
	}

	customFunc.Float[name] = fn

	return nil
}

func RegisterBoolFunc(name string, fn config.BoolCustomFunc) error {
	if _, ok := customFunc.Bool[name]; ok {
		return fail.New(0, "", "API", fail.ErrFuncAlreadyDefined,
			name, "booleans").Error()
	}

	customFunc.Bool[name] = fn

	return nil
}

// Configure passes given options to user configurations.
func Configure(opt *config.Config) {
	userConfig.Configure(opt)
}
