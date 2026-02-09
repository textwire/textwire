package file

import (
	"path/filepath"
	"strings"
)

// joinPaths safely joins 2 paths together treating slashes correctly.
func JoinPaths(path1, path2 string) string {
	return strings.TrimRight(path1, "/") + "/" + strings.TrimLeft(path2, "/")
}

// AppendFileExt adds Textwire file extension to the end of the file if needed.
// It will ignore adding if extension already exist.
func AppendFileExt(path, ext string) string {
	if path == "" || strings.HasSuffix(path, ext) {
		return path
	}
	return path + ext
}

// NameToRelPath turns component and use statement names to relative path
// e.g. layouts/main will be converted to templates/layouts/main.tw
// e.g. components/book will be converted to templates/components/book.tw
func NameToRelPath(name, templDir, ext string) string {
	return JoinPaths(templDir, AppendFileExt(name, ext))
}

// TrimRelPath removes leading / and ./
func TrimRelPath(relPath string) string {
	// Trim ./ from the beginning
	if len(relPath) > 1 && relPath[0] == '.' && relPath[1] == '/' {
		relPath = relPath[2:]
	}

	return strings.TrimLeft(relPath, "/")
}

func ToFullPath(relPath string) (string, error) {
	absPath, err := filepath.Abs(relPath)
	if err != nil {
		return "", err
	}

	return absPath, nil
}
