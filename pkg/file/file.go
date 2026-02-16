package file

import (
	"io"
	"io/fs"
	"os"
	"strings"
	"time"

	"github.com/textwire/textwire/v3/config"
)

// SourceFile holds information about individual Textwire source file,
// including relative and absolute file paths.
type SourceFile struct {
	// Name of the source file, like "components/book" or "layouts/base",
	// or "home". This field can be empty when we evaluate a file outside
	// of using templating system.
	Name string

	// Rel path to the source file that starts from the root of user's project.
	// When using config.TemplateFS, relative path will exclude
	// config.TemplateDir from it to use embeded paths properly.
	Rel string

	// Abs is the absolute path to the source file starting with `/`.
	Abs string

	// ModTime is when the file was last modified.
	ModTime time.Time

	// config is user's configurations that SourceFile needs to access
	// source file extension and location to root of templates.
	config *config.Config
}

func New(name, rel, abs string, c *config.Config) *SourceFile {
	return &SourceFile{
		Name:   strings.Trim(name, "/"),
		Rel:    trimRelPath(rel),
		Abs:    abs,
		config: c,
	}
}

// Content returns the raw content of the file as a string.
func (f *SourceFile) Content() (string, error) {
	var content []byte
	var err error

	if f.config.UsesFS() {
		content, err = fs.ReadFile(f.config.TemplateFS, f.Rel)
	} else {
		content, err = os.ReadFile(f.Abs)
	}

	if err != nil && err != io.EOF {
		return "", err
	}

	return string(content), nil
}
