package textwire

import (
	_ "embed"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/textwire/textwire/v3/pkg/ast"
	"github.com/textwire/textwire/v3/pkg/fail"
	"github.com/textwire/textwire/v3/pkg/file"
	"github.com/textwire/textwire/v3/pkg/lexer"
	"github.com/textwire/textwire/v3/pkg/parser"
)

//go:embed embed/default-error-page.tw
var defaultErrPage string

// errorPage returns HTML that's displayed when an error occurs while
// rendering template.
func errorPage(failure *fail.Error) (string, error) {
	data := map[string]any{
		"path":      failure.Filepath(),
		"line":      failure.Line(),
		"message":   failure.Message(),
		"debugMode": userConf.DebugMode,
	}

	out, err := EvaluateString(defaultErrPage, data)
	if err != nil {
		return "", err
	}

	return out, nil
}

func parseStr(text string) (*ast.Program, []*fail.Error) {
	l := lexer.New(text)
	p := parser.New(l, nil)

	prog := p.ParseProgram()
	if p.HasErrors() {
		return nil, p.Errors()
	}

	return prog, nil
}

// parseFiles parses each Textwire file into AST nodes and returns them.
func parseFiles(files []*file.SourceFile) ([]*ast.Program, *fail.Error) {
	programs := make([]*ast.Program, 0, len(files))
	for _, f := range files {
		prog, failure, parseErr := parseFile(f)
		if parseErr != nil {
			return programs, fail.FromError(parseErr, 0, f.Abs, "template")
		}

		if failure != nil {
			return programs, failure
		}

		programs = append(programs, prog)
	}

	return programs, nil
}

// parseFile parses given file into a ast.Program and returns it.
func parseFile(f *file.SourceFile) (*ast.Program, *fail.Error, error) {
	content, err := f.Content()
	if err != nil {
		return nil, nil, err
	}

	l := lexer.New(content)
	p := parser.New(l, f)
	if p.HasErrors() {
		return nil, p.Errors()[0], nil
	}

	prog := p.ParseProgram()
	prog.AbsPath = f.Abs
	prog.Name = f.Name

	if p.HasErrors() {
		return nil, p.Errors()[0], nil
	}

	return prog, nil, nil
}

// locateFiles recursively finds all Textwire files in the templates directory,
// creates a *file wrapper for each of them, and returns the discovered files.
func locateFiles() ([]*file.SourceFile, error) {
	files := make([]*file.SourceFile, 0, 4)
	err := fs.WalkDir(
		userConf.TemplateFS,
		".",
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() || !strings.Contains(path, userConf.TemplateExt) {
				return nil
			}

			// When using config.TemplateFS to embed templates into binary,
			// we need to exclude config.TemplateDir from path since it
			// already contains it.
			if userConf.UsesFS() {
				path = strings.Replace(path, userConf.TemplateDir, "", 1)
			}

			relPath := file.JoinPaths(userConf.TemplateDir, path)
			absPath, err := filepath.Abs(relPath)
			if err != nil {
				return err
			}

			name := strings.Replace(path, userConf.TemplateExt, "", 1)
			files = append(files, file.New(name, relPath, absPath, userConf))

			return nil
		},
	)

	return files, err
}
