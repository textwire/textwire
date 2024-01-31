package textwire

import (
	"errors"

	"github.com/textwire/textwire/evaluator"
	"github.com/textwire/textwire/object"
)

var config = &Config{
	TemplateDir: "templates",
	TemplateExt: ".textwire.html",
}

// Config is the main configuration for Textwire
type Config struct {
	// TemplateDir is the directory where the Textwire
	// templates are located
	TemplateDir string

	// TemplateExt is the extension of the Textwire
	// template files
	TemplateExt string
}

func New(c *Config) (*Template, error) {
	applyConfig(c)

	paths, err := findTextwireFiles()

	if err != nil {
		return nil, err
	}

	programs, err := parsePrograms(paths)

	return &Template{programs: programs}, err
}

func EvaluateString(inp string, data map[string]interface{}) (string, error) {
	prog, err := parseStr(inp)

	if err != nil {
		return "", err
	}

	env, err := object.EnvFromMap(data)

	if err != nil {
		return "", err
	}

	evaluated := evaluator.Eval(prog, env)

	if evaluated.Is(object.ERROR_OBJ) {
		errMsg := evaluated.(*object.Error).Err.String()
		return "", errors.New(errMsg)
	}

	return evaluated.String(), nil
}
