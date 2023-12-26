package textwire

import (
	"github.com/textwire/textwire/evaluator"
	"github.com/textwire/textwire/lexer"
	"github.com/textwire/textwire/object"
	"github.com/textwire/textwire/parser"
)

func ParseText(text string, vars map[string]interface{}) (string, error) {
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

	if evaluated == nil {
		return "", nil
	}

	return evaluated.String(), nil
}
