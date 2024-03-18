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
	// If you use a different extension other then ".tw.html",
	// you will loose syntax highlighting in VSCode editor
	// if you use the Textwire extension
	TemplateExt string

	// configApplied is a flag to check if the configuration
	// are set or not
	configApplied bool
}

func NewTemplate(c *Config) (*Template, error) {
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
	config.configApplied = false

	prog, errs := parseStr(inp)

	if len(errs) != 0 {
		return "", errs[0].Error()
	}

	env, err := object.EnvFromMap(data)

	if err != nil {
		return "", err.Error()
	}

	ctx := evaluator.NewContext("")
	eval := evaluator.New(ctx)

	evaluated := eval.Eval(prog, env)

	if evaluated.Is(object.ERR_OBJ) {
		return "", evaluated.(*object.Error).Err.Error()
	}

	return evaluated.String(), nil
}

func EvaluateFile(absPath string, data map[string]interface{}) (string, error) {
	config.configApplied = false

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
