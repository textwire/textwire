package textwire

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func getFullPath(filename string, appendExt bool) (string, error) {
	if configApplied {
		filename = conf.TemplateDir + "/" + filename
	}

	if appendExt {
		filename += conf.TemplateExt
	}

	absPath, err := filepath.Abs(filename)

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

	fmt.Printf("-------> %#v\n", conf)
	err := filepath.Walk(conf.TemplateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.Contains(path, conf.TemplateExt) {
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
	name := strings.Replace(path, conf.TemplateDir+"/", "", 1)
	name = strings.Replace(name, conf.TemplateExt, "", 1)
	return name
}
