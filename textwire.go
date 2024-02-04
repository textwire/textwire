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
}

func TemplateEngine(c *Config) (*Template, *fail.Error) {
	applyConfig(c)

	paths, err := findTextwireFiles()

	if err != nil {
		return nil, fail.New(0, "", "template", err.Error())
	}

	programs, errs := parsePrograms(paths)

	if len(errs) != 0 {
		return nil, errs[0]
	}

	return &Template{programs: programs}, nil
}

func EvaluateString(inp string, data map[string]interface{}) (string, []*fail.Error) {
	prog, errs := parseStr(inp)

	if len(errs) != 0 {
		return "", errs
	}

	env, err := object.EnvFromMap(data)

	if err != nil {
		return "", []*fail.Error{err}
	}

	ctx := evaluator.NewContext("")
	eval := evaluator.New(ctx)

	evaluated := eval.Eval(prog, env)

	if evaluated.Is(object.ERR_OBJ) {
		return "", []*fail.Error{
			evaluated.(*object.Error).Err,
		}
	}

	return evaluated.String(), nil
}
