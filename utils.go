package textwire

import (
	"fmt"
	"path/filepath"
)

func getFullPath(fileName string) (string, error) {
	path := fmt.Sprintf("%s/%s%s", config.TemplateDir, fileName, config.TemplateExt)
	absPath, err := filepath.Abs(path)

	if err != nil {
		return "", err
	}

	return absPath, nil
}
