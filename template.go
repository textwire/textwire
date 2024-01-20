package textwire

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/textwire/textwire/ast"
	"github.com/textwire/textwire/evaluator"
	"github.com/textwire/textwire/object"
)

type Template struct {
	prog *ast.Program
}

func (t *Template) Evaluate(vars map[string]interface{}) (object.Object, error) {
	env, err := object.EnvFromMap(vars)

	if err != nil {
		return nil, err
	}

	evaluated := evaluator.Eval(t.prog, env)

	if evaluated.Type() == object.ERROR_OBJ {
		return nil, errors.New(evaluated.String())
	}

	return evaluated, nil
}

func (t *Template) EvaluateResponse(w http.ResponseWriter, vars map[string]interface{}) error {
	evaluated, err := t.Evaluate(vars)

	if err != nil {
		return err
	}

	fmt.Fprint(w, evaluated.String())

	return nil
}
