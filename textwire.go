package textwire

import (
	"github.com/textwire/textwire/evaluator"
	"github.com/textwire/textwire/fail"
	"github.com/textwire/textwire/object"
)

var config = &Config{
	TemplateDir: "templates",
	TemplateExt: ".tw.html",
}

// Config is the main configuration for Textwire
type Config struct {
	// TemplateDir is the directory where the Textwire
	// templates are located. Default is "templates"
	TemplateDir string

	// TemplateExt is the extension of the Textwire
	// template files. Default is ".tw.html"
	TemplateExt string

	// configApplied is a flag to check if the configuration
	// are set or not
	configApplied bool
}

func NewTemplate(c *Config) (*Template, *fail.Error) {
	applyConfig(c)

	paths, err := findTextwireFiles()

	if err != nil {
		return nil, fail.FromError(err, 0, "", "template")
	}

	programs, errs := parsePrograms(paths)

	if len(errs) != 0 {
		return nil, errs[0]
	}

	return &Template{programs: programs}, nil
}

func EvaluateString(inp string, data map[string]interface{}) (string, []*fail.Error) {
	config.configApplied = false

	prog, errs := parseStr(inp)

	if len(errs) != 0 {
		return "", errs
	}

	env, err := object.EnvFromMap(data)

	if err != nil {
		return "", err.ToSlice()
	}

	ctx := evaluator.NewContext("")
	eval := evaluator.New(ctx)

	evaluated := eval.Eval(prog, env)

	if evaluated.Is(object.ERR_OBJ) {
		return "", evaluated.(*object.Error).Err.ToSlice()
	}

	return evaluated.String(), nil
}

func EvaluateFile(absPath string, data map[string]interface{}) (string, []*fail.Error) {
	config.configApplied = false

	_, err := fileContent(absPath)

	if err != nil {
		return "", fail.FromError(err, 0, absPath, "template").ToSlice()
	}

	result, errs := EvaluateString(absPath, data)

	if len(errs) != 0 {
		return "", errs
	}

	return result, nil
}
