package textwire

import (
	"io"
	"io/fs"
	"os"
	"strings"
)

// file holds information about individual Textwire file, including
// relative and absolute file paths.
type file struct {
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
}

func NewFile(name, rel, abs string) *file {
	rel = addTwExtension(rel)
	abs = addTwExtension(abs)

	return &file{
		Name: strings.Trim(name, "/"),
		Rel:  trimRelPath(rel),
		Abs:  abs,
	}
}

// Content returns the content of the file.
func (f *file) Content() (string, error) {
	var content []byte
	var err error

	if userConfig.UsesFS() {
		content, err = fs.ReadFile(userConfig.TemplateFS, f.Rel)
	} else {
		content, err = os.ReadFile(f.Abs)
	}

	if err != nil && err != io.EOF {
		return "", err
	}

	return string(content), nil
}
