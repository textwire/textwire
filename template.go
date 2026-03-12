package textwire

import (
	"fmt"
	"net/http"

	"github.com/textwire/textwire/v3/config"
	"github.com/textwire/textwire/v3/pkg/ast"
	"github.com/textwire/textwire/v3/pkg/evaluator"
	"github.com/textwire/textwire/v3/pkg/fail"
	"github.com/textwire/textwire/v3/pkg/file"
	"github.com/textwire/textwire/v3/pkg/linker"
	"github.com/textwire/textwire/v3/pkg/value"
)

// Template holds all necessary data which it will use when individual
// template files will be evaluated by String() or Response() methods.
type Template struct {
	linker *linker.NodeLinker
}

// NewTemplate returns a new Template instance with parsed Textwire files
// provided by configuration options. The Template instance should be used
// for evaluating Textwire in your handlers.
func NewTemplate(opt *config.Config) (*Template, error) {
	Configure(opt)

	files, err := locateFiles()
	if err != nil {
		return nil, fail.FromError(err, nil, "", "template").Error()
	}

	programs, parseErr := parseFiles(files)
	if parseErr != nil {
		return nil, parseErr.Error()
	}

	ln := linker.New(programs)
	if failure := ln.LinkNodes(); failure != nil {
		return nil, failure.Error()
	}

	tpl := &Template{linker: ln}

	if opt.FileWatcher {
		newFileWatcher(ln).Watch()
	}

	return tpl, nil
}

// String returns final evaluated template result represented as a string.
func (t *Template) String(name string, data map[string]any) (string, *fail.Error) {
	t.linker.RLock()
	linkErr, progs := t.linker.LinkError, t.linker.Programs
	t.linker.RUnlock()

	if linkErr != nil {
		return "", linkErr
	}

	scope, err := value.NewScopeFromMap(data)
	if err != nil {
		return "", err
	}

	name = file.ReplacePathAlias(name, file.PathAliasViews)
	prog := ast.FindProg(name, progs)
	if prog == nil {
		relPath := file.NameToRelPath(name, userConf.TemplateDir, userConf.TemplateExt)
		return "", fail.New(nil, relPath, "template", fail.ErrTemplateNotFound, name)
	}

	e := evaluator.New(customFunc, userConf)
	ctx := evaluator.NewContext(scope, prog.AbsPath)
	evaluated := e.Eval(prog, ctx)
	if evaluated.Is(value.ERR_VAL) {
		return "", evaluated.(*value.Error).Err
	}

	return evaluated.String(), nil
}

// Response evaluates template file with String() method and passing that final
// string to the given http.ResponseWriter.
func (t *Template) Response(w http.ResponseWriter, name string, data map[string]any) error {
	evaluated, failure := t.String(name, data)
	if failure == nil {
		_, err := fmt.Fprint(w, evaluated)
		if err != nil {
			return err
		}

		return nil
	}

	hasErrPage := userConf.ErrorPagePath != ""
	if hasErrPage && !userConf.DebugMode {
		if err := t.responseErrorPage(w); err != nil {
			return err
		}

		return failure.Error()
	}

	errPage, err := errorPage(failure)
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(w, errPage)
	if err != nil {
		return err
	}

	return failure.Error()
}

func (t *Template) responseErrorPage(w http.ResponseWriter) error {
	evaluated, failure := t.String(userConf.ErrorPagePath, nil)
	if failure != nil {
		return failure.Error()
	}

	_, err := fmt.Fprint(w, evaluated)
	if err != nil {
		return err
	}

	return nil
}
