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

// EvaluateString evaluates a given inp string containing Textwire code.
// The function accepts a string template and data to inject into Textwire.
// After evaluation, it returns the processed string and any error encountered.
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

// EvaluateFile evaluates a file containing Textwire code.
//
// The absPath an absolute path to the Textwire file.
// The data is a map of variables you want to inject into the Textwire.
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

// RegisterStrFunc registers a custom function with the given name for the
// string type. You'll be able to use it in your Textwire files.
// e.g. `{{ "Sydney".myFunc() }}`
func RegisterStrFunc(name string, fn config.StrCustomFunc) error {
	if _, ok := customFunc.Str[name]; ok {
		return fail.New(0, "", "API", fail.ErrFuncAlreadyDefined,
			name, "strings").Error()
	}

	customFunc.Str[name] = fn

	return nil
}

// RegisterArrFunc registers a custom function with the given name for the
// array type. You'll be able to use it in your Textwire files.
// e.g. `{{ [1, 2].myFunc() }}`
func RegisterArrFunc(name string, fn config.ArrayCustomFunc) error {
	if _, ok := customFunc.Arr[name]; ok {
		return fail.New(0, "", "API", fail.ErrFuncAlreadyDefined,
			name, "arrays").Error()
	}

	customFunc.Arr[name] = fn

	return nil
}

// RegisterObjFunc registers a custom function with the given name for the
// object type. You'll be able to use it in your Textwire files.
// e.g. `{{ {name: 'Sydney'}.myFunc() }}`
func RegisterObjFunc(name string, fn config.ObjectCustomFunc) error {
	if _, ok := customFunc.Obj[name]; ok {
		return fail.New(0, "", "API", fail.ErrFuncAlreadyDefined,
			name, "objects").Error()
	}

	customFunc.Obj[name] = fn

	return nil
}

// RegisterIntFunc registers a custom function with the given name for the
// integer type. You'll be able to use it in your Textwire files.
// e.g. `{{ 1.myFunc() }}`
func RegisterIntFunc(name string, fn config.IntCustomFunc) error {
	if _, ok := customFunc.Int[name]; ok {
		return fail.New(0, "", "API", fail.ErrFuncAlreadyDefined,
			name, "integers").Error()
	}

	customFunc.Int[name] = fn

	return nil
}

// RegisterFloatFunc registers a custom function with the given name for the
// float type. You'll be able to use it in your Textwire files.
// e.g. `{{ 1.12.myFunc() }}`
func RegisterFloatFunc(name string, fn config.FloatCustomFunc) error {
	if _, ok := customFunc.Float[name]; ok {
		return fail.New(0, "", "API", fail.ErrFuncAlreadyDefined,
			name, "floats").Error()
	}

	customFunc.Float[name] = fn

	return nil
}

// RegisterBoolFunc registers a custom function with the given name for the
// boolean type. You'll be able to use it in your Textwire files.
// e.g. `{{ true.myFunc() }}`
func RegisterBoolFunc(name string, fn config.BoolCustomFunc) error {
	if _, ok := customFunc.Bool[name]; ok {
		return fail.New(0, "", "API", fail.ErrFuncAlreadyDefined,
			name, "booleans").Error()
	}

	customFunc.Bool[name] = fn

	return nil
}

// Configure passes given options to the user configurations.
func Configure(opt *config.Config) {
	userConfig.Configure(opt)
}
