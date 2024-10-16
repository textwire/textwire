package textwire

import (
	"github.com/textwire/textwire/v2/config"
	"github.com/textwire/textwire/v2/evaluator"
	"github.com/textwire/textwire/v2/fail"
	"github.com/textwire/textwire/v2/object"
)

var conf = config.New("templates", ".tw.html")

// usesTemplates is a flag to check if user uses Textwire templates or not
var usesTemplates = false

func NewTemplate(c *config.Config) (*Template, error) {
	applyConfig(c)

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

	ctx := evaluator.NewContext("", conf)
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

func RegisterStrFunc(name string, fn object.BuiltinFunction) error {
	if _, ok := conf.Funcs.Str[name]; ok {
		return fail.New(0, "", "API", fail.ErrFuncAlreadyDefined, name, "strings").Error()
	}

	conf.Funcs.Str[name] = fn

	return nil
}

func RegisterArrFunc(name string, fn object.BuiltinFunction) error {
	if _, ok := conf.Funcs.Arr[name]; ok {
		return fail.New(0, "", "API", fail.ErrFuncAlreadyDefined, name, "arrays").Error()
	}

	conf.Funcs.Arr[name] = fn

	return nil
}

func RegisterIntFunc(name string, fn object.BuiltinFunction) error {
	if _, ok := conf.Funcs.Int[name]; ok {
		return fail.New(0, "", "API", fail.ErrFuncAlreadyDefined, name, "integers").Error()
	}

	conf.Funcs.Int[name] = fn

	return nil
}

func RegisterFloatFunc(name string, fn object.BuiltinFunction) error {
	if _, ok := conf.Funcs.Float[name]; ok {
		return fail.New(0, "", "API", fail.ErrFuncAlreadyDefined, name, "floats").Error()
	}

	conf.Funcs.Float[name] = fn

	return nil
}

func RegisterBoolFunc(name string, fn object.BuiltinFunction) error {
	if _, ok := conf.Funcs.Bool[name]; ok {
		return fail.New(0, "", "API", fail.ErrFuncAlreadyDefined, name, "booleans").Error()
	}

	conf.Funcs.Bool[name] = fn

	return nil
}
