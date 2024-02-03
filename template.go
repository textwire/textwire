package textwire

import (
	"fmt"
	"net/http"

	"github.com/textwire/textwire/ast"
	"github.com/textwire/textwire/evaluator"
	"github.com/textwire/textwire/fail"
	"github.com/textwire/textwire/object"
)

type Template struct {
	programs map[string]*ast.Program
}

func (t *Template) Evaluate(absPath string, data map[string]interface{}) (object.Object, *fail.Error) {
	env, err := object.EnvFromMap(data)

	if err != nil {
		return nil, err
	}

	prog, ok := t.programs[absPath]

	if !ok {
		return nil, fail.New(0, absPath, "template", fail.ErrTemplateNotFound)
	}

	ctx := evaluator.NewContext(absPath)
	eval := evaluator.New(ctx)

	evaluated := eval.Eval(prog, env)

	if evaluated.Is(object.ERR_OBJ) {
		return nil, evaluated.(*object.Error).Err
	}

	return evaluated, nil
}

func (t *Template) View(w http.ResponseWriter, fileName string, data map[string]interface{}) *fail.Error {
	evaluated, err := t.Evaluate(fileName, data)

	if err != nil {
		return err
	}

	fmt.Fprint(w, evaluated.String())

	return nil
}
