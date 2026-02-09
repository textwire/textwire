package file

import (
	"io"
	"io/fs"
	"os"
	"strings"

	"github.com/textwire/textwire/v3/config"
)

// File holds information about individual Textwire File, including
// relative and absolute File paths.
type File struct {
	// Name of the file, like "components/book" or "layouts/base", or "home".
	// This field can be empty when we are evaluating a file outside of using
	// templating system.
	Name string

	// Rel path to the Textwire file that starts from the root of user's
	// project.
	// When using config.TemplateFS, relative path will exclude
	// config.TemplateDir from it to use embeded paths properly.
	Rel string

	// Abs path to the Textwire file starting with `/` and system's root.
	Abs string

	// config is user's config
	config *config.Config
}

func New(name, rel, abs string, c *config.Config) *File {
	rel = AppendFileExt(rel, c.TemplateExt)
	abs = AppendFileExt(abs, c.TemplateExt)

	return &File{
		Name:   strings.Trim(name, "/"),
		Rel:    TrimRelPath(rel),
		Abs:    abs,
		config: c,
	}
}

// Content returns the content of the file.
func (f *File) Content() (string, error) {
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
