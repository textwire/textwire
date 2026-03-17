package textwire

import (
	"github.com/textwire/textwire/v4/config"
	"github.com/textwire/textwire/v4/pkg/evaluator"
	"github.com/textwire/textwire/v4/pkg/fail"
	"github.com/textwire/textwire/v4/pkg/file"
	"github.com/textwire/textwire/v4/pkg/value"
)

var (
	userConf   = config.New("templates", ".tw", "", false)
	customFunc = config.NewFunc()
)

// EvaluateString evaluates a given inp string containing Textwire code.
// The function accepts a string template and data to inject into Textwire.
// After evaluation, it returns the processed string and any error encountered.
func EvaluateString(inp string, data map[string]any) (string, *fail.Error) {
	prog, errs := parseStr(inp)
	if len(errs) != 0 {
		return "", errs[0]
	}

	scope, err := value.NewScopeFromMap(data)
	if err != nil {
		return "", err
	}

	e := evaluator.New(customFunc, nil)
	ctx := evaluator.NewContext(scope, prog.AbsPath)
	evaluated := e.Eval(prog, ctx)
	if evaluated.Is(value.ERR_VAL) {
		return "", evaluated.(*value.Error).Err
	}

	return evaluated.String(), nil
}

// EvaluateFile evaluates a file containing Textwire code.
//
// The absPath an absolute path to the Textwire file.
// The data is a map of variables you want to inject into the Textwire.
func EvaluateFile(absPath string, data map[string]any) (string, *fail.Error) {
	f := file.New("", "", absPath, userConf)
	content, err := f.Content()
	if err != nil {
		return "", fail.FromError(err, nil, absPath, fail.OriginTpl)
	}

	res, failure := EvaluateString(content, data)
	if failure != nil {
		return "", failure
	}

	return res, nil
}

// RegisterStrFunc registers a custom function with the given name for the
// string type. You'll be able to use it in your Textwire files.
// e.g. `{{ "Sydney".myFunc() }}`
func RegisterStrFunc(name string, fn config.StrCustomFunc) *fail.Error {
	if _, ok := customFunc.Str[name]; ok {
		return fail.New(nil, "", fail.OriginAPI, fail.ErrFuncAlreadyDefined, name, "strings")
	}

	customFunc.Str[name] = fn

	return nil
}

// RegisterArrFunc registers a custom function with the given name for the
// array type. You'll be able to use it in your Textwire files.
// e.g. `{{ [1, 2].myFunc() }}`
func RegisterArrFunc(name string, fn config.ArrCustomFunc) *fail.Error {
	if _, ok := customFunc.Arr[name]; ok {
		return fail.New(nil, "", fail.OriginAPI, fail.ErrFuncAlreadyDefined, name, "arrays")
	}

	customFunc.Arr[name] = fn

	return nil
}

// RegisterObjFunc registers a custom function with the given name for the
// object type. You'll be able to use it in your Textwire files.
// e.g. `{{ {name: 'Sydney'}.myFunc() }}`
func RegisterObjFunc(name string, fn config.ObjCustomFunc) *fail.Error {
	if _, ok := customFunc.Obj[name]; ok {
		return fail.New(nil, "", fail.OriginAPI, fail.ErrFuncAlreadyDefined, name, "objects")
	}

	customFunc.Obj[name] = fn

	return nil
}

// RegisterIntFunc registers a custom function with the given name for the
// integer type. You'll be able to use it in your Textwire files.
// e.g. `{{ 1.myFunc() }}`
func RegisterIntFunc(name string, fn config.IntCustomFunc) *fail.Error {
	if _, ok := customFunc.Int[name]; ok {
		return fail.New(nil, "", fail.OriginAPI, fail.ErrFuncAlreadyDefined, name, "integers")
	}

	customFunc.Int[name] = fn

	return nil
}

// RegisterFloatFunc registers a custom function with the given name for the
// float type. You'll be able to use it in your Textwire files.
// e.g. `{{ 1.12.myFunc() }}`
func RegisterFloatFunc(name string, fn config.FloatCustomFunc) *fail.Error {
	if _, ok := customFunc.Float[name]; ok {
		return fail.New(nil, "", fail.OriginAPI, fail.ErrFuncAlreadyDefined, name, "floats")
	}

	customFunc.Float[name] = fn

	return nil
}

// RegisterBoolFunc registers a custom function with the given name for the
// boolean type. You'll be able to use it in your Textwire files.
// e.g. `{{ true.myFunc() }}`
func RegisterBoolFunc(name string, fn config.BoolCustomFunc) *fail.Error {
	if _, ok := customFunc.Bool[name]; ok {
		return fail.New(nil, "", fail.OriginAPI, fail.ErrFuncAlreadyDefined, name, "booleans")
	}

	customFunc.Bool[name] = fn

	return nil
}

// Configure passes given options to the user configurations.
func Configure(opt *config.Config) {
	userConf.Configure(opt)
}
