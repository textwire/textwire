package textwire

import (
	"fmt"
	"net/http"

	"github.com/textwire/textwire/v2/ast"
	"github.com/textwire/textwire/v2/ctx"
	"github.com/textwire/textwire/v2/evaluator"
	"github.com/textwire/textwire/v2/fail"
	"github.com/textwire/textwire/v2/object"
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
		return "", fail.New(0, filename, "template", "%s", err.Error())
	}

	prog, ok := t.programs[filename]

	if !ok {
		return "", fail.New(0, absPath, "template", fail.ErrTemplateNotFound)
	}

	ctx := ctx.NewContext(absPath, customFunc, userConfig)
	eval := evaluator.New(ctx)

	evaluated := eval.Eval(prog, env)

	if evaluated.Is(object.ERR_OBJ) {
		return "", evaluated.(*object.Error).Err
	}

	return evaluated.String(), nil
}

func (t *Template) Response(w http.ResponseWriter, filename string, data map[string]interface{}) error {
	evaluated, failErr := t.String(filename, data)

	if failErr == nil {
		fmt.Fprint(w, evaluated)
		return nil
	}

	if userConfig.ErrorPagePath == "" {
		fmt.Fprint(w, errorPage(failErr))
		return failErr.Error()
	}

	if err := t.responseErrorPage(w); err != nil {
		return err
	}

	return failErr.Error()
}

func (t *Template) responseErrorPage(w http.ResponseWriter) error {
	evaluated, err := t.String(userConfig.ErrorPagePath, nil)

	if err != nil {
		return err.Error()
	}

	fmt.Fprint(w, evaluated)

	return nil
}
