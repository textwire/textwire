package textwire

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/textwire/textwire/ast"
	"github.com/textwire/textwire/evaluator"
	"github.com/textwire/textwire/object"
)

type View struct {
	program *ast.Program
}

func (v *View) Evaluate(w http.ResponseWriter, vars map[string]interface{}) error {
	env, err := object.EnvFromMap(vars)

	if err != nil {
		return err
	}

	evaluated := evaluator.Eval(v.program, env)

	if evaluated.Type() == object.ERROR_OBJ {
		return errors.New(evaluated.String())
	}

	fmt.Fprint(w, evaluated.String())

	return nil
}
