package textwire

import (
	"github.com/textwire/textwire/v3/config"
	"github.com/textwire/textwire/v3/pkg/evaluator"
	"github.com/textwire/textwire/v3/pkg/fail"
	"github.com/textwire/textwire/v3/pkg/file"
	"github.com/textwire/textwire/v3/pkg/object"
)

var (
	userConf   = config.New("templates", ".tw", "", false)
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

	scope, err := object.NewScopeFromMap(data)
	if err != nil {
		return "", err.Error()
	}

	e := evaluator.New(customFunc, nil)
	ctx := evaluator.NewContext(scope, prog.AbsPath)
	evaluated := e.Eval(prog, ctx)
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
	f := file.New("", "", absPath, userConf)
	content, err := f.Content()
	if err != nil {
		return "", fail.FromError(err, 0, absPath, "template").Error()
	}

	res, err := EvaluateString(content, data)
	if err != nil {
		return "", err
	}

	return res, nil
}

// RegisterStrFunc registers a custom function with the given name for the
// string type. You'll be able to use it in your Textwire files.
// e.g. `{{ "Sydney".myFunc() }}`
func RegisterStrFunc(name string, fn config.StringCustomFunc) error {
	if _, ok := customFunc.String[name]; ok {
		return fail.New(0, "", "API", fail.ErrFuncAlreadyDefined, name, "strings").Error()
	}

	customFunc.String[name] = fn

	return nil
}

// RegisterArrFunc registers a custom function with the given name for the
// array type. You'll be able to use it in your Textwire files.
// e.g. `{{ [1, 2].myFunc() }}`
func RegisterArrFunc(name string, fn config.ArrayCustomFunc) error {
	if _, ok := customFunc.Array[name]; ok {
		return fail.New(0, "", "API", fail.ErrFuncAlreadyDefined, name, "arrays").Error()
	}

	customFunc.Array[name] = fn

	return nil
}

// RegisterObjFunc registers a custom function with the given name for the
// object type. You'll be able to use it in your Textwire files.
// e.g. `{{ {name: 'Sydney'}.myFunc() }}`
func RegisterObjFunc(name string, fn config.MapCustomFunc) error {
	if _, ok := customFunc.Map[name]; ok {
		return fail.New(0, "", "API", fail.ErrFuncAlreadyDefined, name, "objects").Error()
	}

	customFunc.Map[name] = fn

	return nil
}

// RegisterIntFunc registers a custom function with the given name for the
// integer type. You'll be able to use it in your Textwire files.
// e.g. `{{ 1.myFunc() }}`
func RegisterIntFunc(name string, fn config.IntegerCustomFunc) error {
	if _, ok := customFunc.Integer[name]; ok {
		return fail.New(0, "", "API", fail.ErrFuncAlreadyDefined, name, "integers").Error()
	}

	customFunc.Integer[name] = fn

	return nil
}

// RegisterFloatFunc registers a custom function with the given name for the
// float type. You'll be able to use it in your Textwire files.
// e.g. `{{ 1.12.myFunc() }}`
func RegisterFloatFunc(name string, fn config.FloatCustomFunc) error {
	if _, ok := customFunc.Float[name]; ok {
		return fail.New(0, "", "API", fail.ErrFuncAlreadyDefined, name, "floats").Error()
	}

	customFunc.Float[name] = fn

	return nil
}

// RegisterBoolFunc registers a custom function with the given name for the
// boolean type. You'll be able to use it in your Textwire files.
// e.g. `{{ true.myFunc() }}`
func RegisterBoolFunc(name string, fn config.BooleanCustomFunc) error {
	if _, ok := customFunc.Boolean[name]; ok {
		return fail.New(0, "", "API", fail.ErrFuncAlreadyDefined, name, "booleans").Error()
	}

	customFunc.Boolean[name] = fn

	return nil
}

// Configure passes given options to the user configurations.
func Configure(opt *config.Config) {
	userConf.Configure(opt)
}
