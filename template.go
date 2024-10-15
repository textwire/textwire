package textwire

import (
	"fmt"
	"net/http"

	ast "github.com/textwire/textwire/v2/ast"
	evaluator "github.com/textwire/textwire/v2/evaluator"
	fail "github.com/textwire/textwire/v2/fail"
	object "github.com/textwire/textwire/v2/object"
)

type Template struct {
	programs map[string]*ast.Program
}

func (t *Template) String(filename string, data map[string]interface{}) (string, *fail.Error) {
	env, envErr := object.EnvFromMap(data)

	if envErr != nil {
		return "", envErr
	}

	absPath, err := getFullPath(filename, true)

	if err != nil {
		return "", fail.New(0, filename, "template", err.Error())
	}

	prog, ok := t.programs[filename]

	if !ok {
		return "", fail.New(0, absPath, "template", fail.ErrTemplateNotFound)
	}

	ctx := evaluator.NewContext(absPath)
	eval := evaluator.New(ctx)

	evaluated := eval.Eval(prog, env)

	if evaluated.Is(object.ERR_OBJ) {
		return "", evaluated.(*object.Error).Err
	}

	return evaluated.String(), nil
}

func (t *Template) Response(w http.ResponseWriter, filename string, data map[string]interface{}) error {
	evaluated, err := t.String(filename, data)

	if err != nil {
		return err.Error()
	}

	fmt.Fprint(w, evaluated)

	return nil
}
