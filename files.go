package textwire

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func getFullPath(filename string) (string, error) {
	if usingTemplates {
		filename = joinPaths(userConfig.TemplateDir, filename)
	}

	addTwExtension(filename)

	absPath, err := filepath.Abs(filename)
	if err != nil {
		return "", err
	}

	return absPath, nil
}

func joinPaths(path1, path2 string) string {
	return strings.TrimRight(path1, "/") + "/" + strings.TrimLeft(path2, "/")
}

func fileContent(path string) (string, error) {
	var content []byte
	var err error

	isAbsPath := path[0] == '/'

	if isAbsPath {
		content, err = os.ReadFile(path)
	} else {
		content, err = fs.ReadFile(userConfig.TemplateFS, path)
	}

	if err != nil && err != io.EOF {
		return "", err
	}

	return string(content), nil
}

// findTextwireFiles recursively finds all files in the template
// directory and its nested subdirectories
func findTextwireFiles() (map[string]string, error) {
	var result = map[string]string{}

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

			absPath, err := filepath.Abs(fmt.Sprintf("%s/%s", userConfig.TemplateDir, path))
			if err != nil {
				return err
			}

			result[path] = absPath

			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	return result, nil
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
	return userConfig.TemplateDir + "/" + addTwExtension(name)
}
