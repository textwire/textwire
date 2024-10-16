package textwire

import (
	"github.com/textwire/textwire/v2/evaluator"
	"github.com/textwire/textwire/v2/fail"
	"github.com/textwire/textwire/v2/object"
	"github.com/textwire/textwire/v2/option"
)

var conf = option.New("templates", ".tw.html")
var customFunc = option.NewFunc()

// usesTemplates is a flag to check if user uses Textwire templates or not
var usesTemplates = false

func NewTemplate(opt *option.Option) (*Template, error) {
	applyOptions(opt)

	paths, err := findTextwireFiles()

	if err != nil {
		return nil, fail.FromError(err, 0, "", "template").Error()
	}

	programs, parseErr := parsePrograms(paths)

	if parseErr != nil {
		return nil, parseErr.Error()
	}

	return &Template{programs: programs}, nil
}

func EvaluateString(inp string, data map[string]interface{}) (string, error) {
	usesTemplates = false

	prog, errs := parseStr(inp)

	if len(errs) != 0 {
		return "", errs[0].Error()
	}

	env, err := object.EnvFromMap(data)

	if err != nil {
		return "", err.Error()
	}

	ctx := evaluator.NewContext("", customFunc)
	eval := evaluator.New(ctx)

	evaluated := eval.Eval(prog, env)

	if evaluated.Is(object.ERR_OBJ) {
		return "", evaluated.(*object.Error).Err.Error()
	}

	return evaluated.String(), nil
}

func EvaluateFile(absPath string, data map[string]interface{}) (string, error) {
	usesTemplates = false

	_, err := fileContent(absPath)

	if err != nil {
		return "", fail.FromError(err, 0, absPath, "template").Error()
	}

	result, err := EvaluateString(absPath, data)

	if err != nil {
		return "", err
	}

	return result, nil
}

func RegisterStrFunc(name string, fn option.StrCustomFunc) error {
	if _, ok := customFunc.Str[name]; ok {
		return fail.New(0, "", "API", fail.ErrFuncAlreadyDefined, name, "strings").Error()
	}

	customFunc.Str[name] = fn

	return nil
}

func RegisterArrFunc(name string, fn option.ArrayCustomFunc) error {
	if _, ok := customFunc.Arr[name]; ok {
		return fail.New(0, "", "API", fail.ErrFuncAlreadyDefined, name, "arrays").Error()
	}

	customFunc.Arr[name] = fn

	return nil
}

func RegisterIntFunc(name string, fn option.IntCustomFunc) error {
	if _, ok := customFunc.Int[name]; ok {
		return fail.New(0, "", "API", fail.ErrFuncAlreadyDefined, name, "integers").Error()
	}

	customFunc.Int[name] = fn

	return nil
}

func RegisterFloatFunc(name string, fn option.FloatCustomFunc) error {
	if _, ok := customFunc.Float[name]; ok {
		return fail.New(0, "", "API", fail.ErrFuncAlreadyDefined, name, "floats").Error()
	}

	customFunc.Float[name] = fn

	return nil
}

func RegisterBoolFunc(name string, fn option.BoolCustomFunc) error {
	if _, ok := customFunc.Bool[name]; ok {
		return fail.New(0, "", "API", fail.ErrFuncAlreadyDefined, name, "booleans").Error()
	}

	customFunc.Bool[name] = fn

	return nil
}
