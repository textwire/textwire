package textwire

import (
	_ "embed"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

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

// findTwFiles recursively finds all files in the templates directory,
// and creates a *twFile wrapper for each of these files.
func findTwFiles() ([]*textwireFile, error) {
	twPaths := make([]*textwireFile, 0, 4) // 4 is an approximate number

	err := fs.WalkDir(
		userConfig.TemplateFS,
		".",
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() || !strings.Contains(path, userConfig.TemplateExt) {
				return nil
			}

			// When using config.TemplateFS to embed templates into binary,
			// we need to exclude config.TemplateDir from path since it
			// already contains it.
			if userConfig.UsesFS() {
				path = strings.Replace(path, userConfig.TemplateDir, "", 1)
			}

			relPath := joinPaths(userConfig.TemplateDir, path)
			absPath, err := filepath.Abs(relPath)
			if err != nil {
				return err
			}

			twPaths = append(twPaths, NewTextwireFile(relPath, absPath))

			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	return twPaths, nil
}

// fileContent returns the content of the provided file path.
func fileContent(twFile *textwireFile) (string, error) {
	var content []byte
	var err error

	if userConfig.UsesFS() {
		content, err = fs.ReadFile(userConfig.TemplateFS, twFile.Rel)
	} else {
		content, err = os.ReadFile(twFile.Abs)
	}

	if err != nil && err != io.EOF {
		return "", err
	}

	return string(content), nil
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
