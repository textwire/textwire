package textwire

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/textwire/textwire/v3/ast"
	"github.com/textwire/textwire/v3/fail"
	"github.com/textwire/textwire/v3/lexer"
	"github.com/textwire/textwire/v3/parser"
)

// twFile information about individual Textwire File, like path and name.
type twFile struct {
	// Rel path to the Textwire file that starts from the root of user's
	// project.
	// When using config.TemplateFS, relative path will exclude
	// config.TemplateDir from it to use embeded paths properly.
	Rel string

	// Abs path to the Textwire file starting with `/` and system's root.
	Abs string

	// Prog is parsed AST for this Textwire file.
	Prog *ast.Program
}

func NewTwFile(rel, abs string) *twFile {
	rel = addTwExtension(rel)
	abs = addTwExtension(abs)

	return &twFile{
		Rel: trimRelPath(rel),
		Abs: abs,
	}
}

// parseProgram parses twFile.Program and returns errors.
func (tf *twFile) parseProgram() (*fail.Error, error) {
	content, err := fileContent(tf)
	if err != nil {
		return nil, err
	}

	lex := lexer.New(content)
	pars := parser.New(lex, tf.Abs)
	tf.Prog = pars.ParseProgram()

	if len(pars.Errors()) != 0 {
		return pars.Errors()[0], nil
	}

	tf.Prog.Filepath = tf.Abs

	return nil, nil
}

// findTwFiles recursively finds all files in the templates directory,
// and creates a *twFile wrapper for each of these files.
func findTwFiles() ([]*twFile, error) {
	twPaths := make([]*twFile, 0, 4) // 4 is an approximate number

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

			twPaths = append(twPaths, NewTwFile(relPath, absPath))

			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	return twPaths, nil
}

// fileContent returns the content of the provided file path.
func fileContent(twFile *twFile) (string, error) {
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
