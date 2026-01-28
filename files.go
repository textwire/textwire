package textwire

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

func getFullPath(filename string, appendExt bool) (string, error) {
	if usingTemplates {
		filename = joinPaths(userConfig.TemplateDir, filename)
	}

	if appendExt {
		filename += userConfig.TemplateExt
	}

	absPath, err := filepath.Abs(filename)
	if err != nil {
		return "", err
	}

	return absPath, nil
}

func joinPaths(path1, path2 string) string {
	return strings.TrimRight(path1, "/") + "/" + strings.TrimLeft(path2, "/")
}

func fileContent(absPath string) (string, error) {
	content, err := os.ReadFile(absPath)

	if err != nil && err != io.EOF {
		return "", err
	}

	return string(content), nil
}

// findTextwireFiles recursively finds all files in the template
// directory and its nested subdirectories
func findTextwireFiles() (map[string]string, error) {
	var result = map[string]string{}

	err := filepath.Walk(userConfig.TemplateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.Contains(path, userConfig.TemplateExt) {
			return nil
		}

		absPath, err := filepath.Abs(path)

		if err != nil {
			return err
		}

		result[nameFromPath(path)] = absPath

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func nameFromPath(path string) string {
	name := strings.Replace(path, userConfig.TemplateDir+"/", "", 1)
	name = strings.Replace(name, userConfig.TemplateExt, "", 1)
	return name
}
