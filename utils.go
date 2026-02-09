package textwire

import (
	_ "embed"
	"path/filepath"
	"strings"

	"github.com/textwire/textwire/v3/ast"
	"github.com/textwire/textwire/v3/fail"
	"github.com/textwire/textwire/v3/lexer"
	"github.com/textwire/textwire/v3/parser"
)

//go:embed textwire/default-error-page.tw
var defaultErrPage string

// errorPage returns HTML that's displayed when an error occurs while
// rendering template.
func errorPage(failure *fail.Error) (string, error) {
	data := map[string]any{
		"path":      failure.Filepath(),
		"line":      failure.Line(),
		"message":   failure.Message(),
		"debugMode": userConfig.DebugMode,
	}

	result, err := EvaluateString(defaultErrPage, data)
	if err != nil {
		return "", err
	}

	return result, nil
}

func getFullPath(relPath string) (string, error) {
	absPath, err := filepath.Abs(relPath)
	if err != nil {
		return "", err
	}

	return absPath, nil
}

// joinPaths safely joins 2 paths together treating slashes correctly.
func joinPaths(path1, path2 string) string {
	return strings.TrimRight(path1, "/") + "/" + strings.TrimLeft(path2, "/")
}

// trimRelPath removes leading / and ./
func trimRelPath(relPath string) string {
	// Trim ./ from the beginning
	if len(relPath) > 1 && relPath[0] == '.' && relPath[1] == '/' {
		relPath = relPath[2:]
	}

	return strings.TrimLeft(relPath, "/")
}

// addTwExtension adds Textwire file extension to the end of the file if needed.
// It will ignore adding if extension already exist.
func addTwExtension(path string) string {
	if path == "" || strings.HasSuffix(path, userConfig.TemplateExt) {
		return path
	}
	return path + userConfig.TemplateExt
}

// nameToRelPath turns component and use statement names to relative path
// e.g. layouts/main will be converted to templates/layouts/main.tw
// e.g. components/book will be converted to templates/components/book.tw
func nameToRelPath(name string) string {
	return joinPaths(userConfig.TemplateDir, addTwExtension(name))
}

func parseStr(text string) (*ast.Program, []*fail.Error) {
	l := lexer.New(text)
	p := parser.New(l, "", "")

	prog := p.ParseProgram()
	if p.HasErrors() {
		return nil, p.Errors()
	}

	return prog, nil
}
