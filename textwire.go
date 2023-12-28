package textwire

import (
	"errors"
	"os"

	"github.com/textwire/textwire/evaluator"
	"github.com/textwire/textwire/lexer"
	"github.com/textwire/textwire/object"
	"github.com/textwire/textwire/parser"
)

// ParseStr parses a Textwire string and returns the result as a string
func ParseStr(text string, vars map[string]interface{}) (string, error) {
	lex := lexer.New(text)
	pars := parser.New(lex)

	program := pars.ParseProgram()

	if len(pars.Errors()) != 0 {
		return "", pars.Errors()[0]
	}

	env, err := object.EnvFromMap(vars)

	if err != nil {
		return "", err
	}

	evaluated := evaluator.Eval(program, env)

	if evaluated.Type() == object.ERROR_OBJ {
		return "", errors.New(evaluated.String())
	}

	return evaluated.String(), nil
}

// ParseFile parses a Textwire file and returns the result as a string
func ParseFile(filePath string, vars map[string]interface{}) (string, error) {
	content, err := os.ReadFile(filePath)

	if err != nil {
		return "", err
	}

	strContent := string(content)

	return ParseStr(strContent, vars)
}
