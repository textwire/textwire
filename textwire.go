package textwire

import (
	"strings"

	"github.com/textwire/textwire/v2/config"
	"github.com/textwire/textwire/v2/ctx"
	"github.com/textwire/textwire/v2/evaluator"
	"github.com/textwire/textwire/v2/fail"
	"github.com/textwire/textwire/v2/object"
)

var userConfig = config.New("templates", ".tw.html", "", false)
var customFunc = config.NewFunc()

// usesTemplates is a flag to check if user uses Textwire templates or not
var usesTemplates = false

func NewTemplate(opt *config.Config) (*Template, error) {
	Configure(opt)

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

	ctx := ctx.NewContext("", customFunc, userConfig)
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

func RegisterStrFunc(name string, fn config.StrCustomFunc) error {
	if _, ok := customFunc.Str[name]; ok {
		return fail.New(0, "", "API", fail.ErrFuncAlreadyDefined, name, "strings").Error()
	}

	customFunc.Str[name] = fn

	return nil
}

func RegisterArrFunc(name string, fn config.ArrayCustomFunc) error {
	if _, ok := customFunc.Arr[name]; ok {
		return fail.New(0, "", "API", fail.ErrFuncAlreadyDefined, name, "arrays").Error()
	}

	customFunc.Arr[name] = fn

	return nil
}

func RegisterIntFunc(name string, fn config.IntCustomFunc) error {
	if _, ok := customFunc.Int[name]; ok {
		return fail.New(0, "", "API", fail.ErrFuncAlreadyDefined, name, "integers").Error()
	}

	customFunc.Int[name] = fn

	return nil
}

func RegisterFloatFunc(name string, fn config.FloatCustomFunc) error {
	if _, ok := customFunc.Float[name]; ok {
		return fail.New(0, "", "API", fail.ErrFuncAlreadyDefined, name, "floats").Error()
	}

	customFunc.Float[name] = fn

	return nil
}

func RegisterBoolFunc(name string, fn config.BoolCustomFunc) error {
	if _, ok := customFunc.Bool[name]; ok {
		return fail.New(0, "", "API", fail.ErrFuncAlreadyDefined, name, "booleans").Error()
	}

	customFunc.Bool[name] = fn

	return nil
}

func Configure(opt *config.Config) {
	usesTemplates = true

	if opt == nil {
		return
	}

	if opt.TemplateDir != "" {
		userConfig.TemplateDir = strings.Trim(opt.TemplateDir, "/")
	}

	if opt.TemplateExt != "" {
		userConfig.TemplateExt = opt.TemplateExt
	}

	userConfig.DebugMode = opt.DebugMode
}
