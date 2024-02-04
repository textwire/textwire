package textwire

import (
	"os"
	"path/filepath"
	"strings"
)

func getFullPath(filename string) (string, error) {
	path := config.TemplateDir + "/" + filename + config.TemplateExt
	absPath, err := filepath.Abs(path)

	if err != nil {
		return "", err
	}

	return absPath, nil
}

func fileContent(absPath string) (string, error) {
	content, err := os.ReadFile(absPath)

	if err != nil {
		return "", err
	}

	return string(content), nil
}

// findTextwireFiles recursively finds all files in the template
// directory and its nested subdirectories
func findTextwireFiles() (map[string]string, error) {
	var result = map[string]string{}

	err := filepath.Walk(config.TemplateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.Contains(path, config.TemplateExt) {
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
	name := strings.Replace(path, config.TemplateDir+"/", "", 1)
	name = strings.Replace(name, config.TemplateExt, "", 1)
	return name
}
