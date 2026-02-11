package textwire

import (
	_ "embed"

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
		"debugMode": userConf.DebugMode,
	}

	out, err := EvaluateString(defaultErrPage, data)
	if err != nil {
		return "", err
	}

	return result, nil
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
