package textwire

import (
	"fmt"
	"net/http"
	"strings"

	"codeberg.org/textwire/textwire/v5/config"
	"codeberg.org/textwire/textwire/v5/pkg/ast"
	"codeberg.org/textwire/textwire/v5/pkg/evaluator"
	"codeberg.org/textwire/textwire/v5/pkg/fail"
	"codeberg.org/textwire/textwire/v5/pkg/file"
	"codeberg.org/textwire/textwire/v5/pkg/linker"
	"codeberg.org/textwire/textwire/v5/pkg/value"
)

const (
	tplNotFoundPrefix = "template file '"
	tplNotFoundSuffix = "' not found"
)

// Template holds all necessary data which it will use when individual
// template files will be evaluated by String() or Response() methods.
type Template struct {
	linker *linker.NodeLinker
}

// NewTemplate returns a new Template instance with parsed Textwire files
// provided by configuration options. The Template instance should be used
// for evaluating Textwire in your handlers.
func NewTemplate(opt *config.Config) (*Template, *fail.Error) {
	Configure(opt)

	files, err := locateFiles()
	if err != nil {
		return nil, fail.FromError(err, nil, "", fail.OriginTpl)
	}

	programs, parseErr := parseFiles(files)
	if parseErr != nil {
		return nil, parseErr
	}

	ln := linker.New(programs)
	if failure := ln.LinkNodes(); failure != nil {
		return nil, failure
	}

	tpl := &Template{linker: ln}

	if userConf.FileWatcher {
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
		return "", fail.New(nil, relPath, fail.OriginTpl, fail.ErrTemplateNotFound, name)
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
func (t *Template) Response(w http.ResponseWriter, name string, data map[string]any) *fail.Error {
	evaluated, failure := t.String(name, data)
	if failure == nil {
		return writeBody(w, name, evaluated)
	}

	// Rendering the error response itself failed: report that error instead.
	if err := t.writeFailure(w, name, failure); err != nil {
		return err
	}

	return failure
}

// isTemplateNotFound reports whether the failure is caused by a missing
// template file.
func isTemplateNotFound(failure *fail.Error) bool {
	if failure == nil || failure.Origin() != fail.OriginTpl {
		return false
	}

	msg := failure.Message()
	return strings.HasPrefix(msg, tplNotFoundPrefix) &&
		strings.HasSuffix(msg, tplNotFoundSuffix)
}

// httpStatusFor maps a template failure to the HTTP status code to respond with.
func httpStatusFor(failure *fail.Error) int {
	if isTemplateNotFound(failure) {
		return http.StatusNotFound
	}

	return http.StatusInternalServerError
}

// writeBody writes the given content to w, wrapping any write error.
func writeBody(w http.ResponseWriter, name, content string) *fail.Error {
	if _, err := fmt.Fprint(w, content); err != nil {
		return fail.FromError(err, nil, name, fail.OriginTpl)
	}

	return nil
}

// writeFailure writes the error response (custom error page or the built-in
// one) for the given failure, using the appropriate HTTP status.
func (t *Template) writeFailure(
	w http.ResponseWriter,
	name string,
	failure *fail.Error,
) *fail.Error {
	status := httpStatusFor(failure)

	if userConf.ErrorPagePath != "" && !userConf.DebugMode {
		w.WriteHeader(status)
		return t.responseCustomErrorPage(w)
	}

	page, err := errorPage(failure)
	if err != nil {
		return err
	}

	w.WriteHeader(status)
	return writeBody(w, name, page)
}

func (t *Template) responseCustomErrorPage(w http.ResponseWriter) *fail.Error {
	evaluated, failure := t.String(userConf.ErrorPagePath, nil)
	if failure != nil {
		return failure
	}

	_, err := fmt.Fprint(w, evaluated)
	if err != nil {
		return fail.FromError(err, nil, userConf.ErrorPagePath, fail.OriginTpl)
	}

	return nil
}
