package textwire

import (
	"fmt"
	"net/http"

	"github.com/textwire/textwire/v3/ast"
	"github.com/textwire/textwire/v3/evaluator"
	"github.com/textwire/textwire/v3/fail"
	"github.com/textwire/textwire/v3/object"
)

type Template struct {
	programs map[string]*ast.Program
}

func (t *Template) String(filename string, data map[string]any) (string, *fail.Error) {
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

	prog.Filepath = absPath
	eval := evaluator.New(customFunc, userConfig)

	evaluated := eval.Eval(prog, env, prog.Filepath)
	if evaluated.Is(object.ERR_OBJ) {
		return "", evaluated.(*object.Error).Err
	}

	return evaluated.String(), nil
}

func (t *Template) Response(w http.ResponseWriter, filename string, data map[string]any) error {
	evaluated, failure := t.String(filename, data)
	if failure == nil {
		_, err := fmt.Fprint(w, evaluated)
		if err != nil {
			return err
		}

		return nil
	}

	hasErrorPage := userConfig.ErrorPagePath != ""
	if hasErrorPage && !userConfig.DebugMode {
		if err := t.responseErrorPage(w); err != nil {
			return err
		}

		return failure.Error()
	}

	out, err := errorPage(failure)
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(w, out)
	if err != nil {
		return err
	}

	return failure.Error()
}

func (t *Template) responseErrorPage(w http.ResponseWriter) error {
	evaluated, failure := t.String(userConfig.ErrorPagePath, nil)
	if failure != nil {
		return failure.Error()
	}

	_, err := fmt.Fprint(w, evaluated)
	if err != nil {
		return err
	}

	return nil
}
