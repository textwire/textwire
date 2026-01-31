package textwire

import (
	_ "embed"

	"github.com/textwire/textwire/v3/fail"
)

//go:embed textwire/default-error-page.tw
var defaultErrorPage string

// errorPage returns HTML that's displayed when an error
// occurs while rendering a template
func errorPage(failure *fail.Error) (string, error) {
	data := map[string]any{
		"path":      failure.Filepath(),
		"line":      failure.Line(),
		"message":   failure.Message(),
		"debugMode": userConfig.DebugMode,
	}

	result, err := EvaluateString(defaultErrorPage, data)
	if err != nil {
		return "", err
	}

	return result, nil
}
