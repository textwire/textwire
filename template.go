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
	programs map[string]*ast.Program
}

func (t *Template) Evaluate(fileName string, data map[string]interface{}) (object.Object, error) {
	env, err := object.EnvFromMap(data)

	if err != nil {
		return nil, err
	}

	prog, ok := t.programs[fileName]

	if !ok {
		return nil, fmt.Errorf("template \"%s\" not found", fileName)
	}

	evaluated := evaluator.Eval(prog, env)

	if evaluated.Is(object.ERR_OBJ) {
		return nil, errors.New(evaluated.String())
	}

	return evaluated, nil
}

func (t *Template) View(w http.ResponseWriter, fileName string, data map[string]interface{}) error {
	evaluated, err := t.Evaluate(fileName, data)

	if err != nil {
		return err
	}

	fmt.Fprint(w, evaluated.String())

	return nil
}
