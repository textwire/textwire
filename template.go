package textwire

import (
	"fmt"
	"net/http"

	"github.com/textwire/textwire/v3/config"
	"github.com/textwire/textwire/v3/evaluator"
	"github.com/textwire/textwire/v3/fail"
	"github.com/textwire/textwire/v3/object"
)

type Template struct {
	sourceBundle *SourceBundle
}

// NewTemplate returns a new Template instance with parsed Textwire files
// provided by configuration options. The Template instance should be used
// for evaluating Textwire in your handlers.
func NewTemplate(opt *config.Config) (*Template, error) {
	Configure(opt)

	sb := NewSourceBundle()

	if err := sb.FindFiles(); err != nil {
		return nil, fail.FromError(err, 0, "", "template").Error()
	}

	if err := sb.ParseFiles(); err != nil {
		return nil, err.Error()
	}

	if err := sb.AddAttachments(); err != nil {
		return nil, err.Error()
	}

	return &Template{sourceBundle: sb}, nil
}

func (t *Template) String(name string, data map[string]any) (string, *fail.Error) {
	scope, err := object.NewScopeFromMap(data)
	if err != nil {
		return "", err
	}

	prog := t.sourceBundle.FindProg(name)
	if prog == nil {
		return "", fail.New(0, nameToRelPath(name), "template", fail.ErrTemplateNotFound, name)
	}

	e := evaluator.New(customFunc, userConfig)
	evaluated := e.Eval(prog, scope, prog.AbsPath)
	if evaluated.Is(object.ERR_OBJ) {
		return "", evaluated.(*object.Error).Err
	}

	return evaluated.String(), nil
}

func (t *Template) Response(w http.ResponseWriter, name string, data map[string]any) error {
	evaluated, failure := t.String(name, data)
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
