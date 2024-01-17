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
	program *ast.Program
}

func (t *Template) EvaluateResponse(w http.ResponseWriter, vars map[string]interface{}) error {
	env, err := object.EnvFromMap(vars)

	if err != nil {
		return err
	}

	evaluated := evaluator.Eval(t.program, env)

	if evaluated.Type() == object.ERROR_OBJ {
		return errors.New(evaluated.String())
	}

	fmt.Fprint(w, evaluated.String())

	return nil
}
