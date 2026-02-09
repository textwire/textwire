package textwire

import (
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/textwire/textwire/v3/ast"
	"github.com/textwire/textwire/v3/fail"
	"github.com/textwire/textwire/v3/lexer"
	"github.com/textwire/textwire/v3/parser"
)

// SourceBundler is the main struct to handle parsing and evaluation of
// Textwire code.
type SourceBundler struct {
	files []*file
}

func NewSourceBundle() *SourceBundler {
	return &SourceBundler{
		files: make([]*file, 0, 4),
	}
}

// ParseFiles parses each Textwire file into AST nodes and returns them.
func (sb *SourceBundler) ParseFiles() ([]*ast.Program, *fail.Error) {
	programs := make([]*ast.Program, 0, 4)
	for _, f := range sb.files {
		prog, failure, parseErr := sb.parseFile(f)
		if parseErr != nil {
			return programs, fail.FromError(parseErr, 0, f.Abs, "template")
		}

		if failure != nil {
			return programs, failure
		}

		programs = append(programs, prog)
	}

	return programs, nil
}

// FindFiles recursively finds all Textwire files in the templates directory,
// and creates a *file wrapper for each of these files.
func (sb *SourceBundler) FindFiles() error {
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

			name := strings.Replace(path, userConfig.TemplateExt, "", 1)
			sb.files = append(sb.files, NewFile(name, relPath, absPath))

			return nil
		},
	)

	return err
}

// parseFile parses given file into a ast.Program and returns it.
func (sb *SourceBundler) parseFile(f *file) (*ast.Program, *fail.Error, error) {
	content, err := f.Content()
	if err != nil {
		return nil, nil, err
	}

	l := lexer.New(content)
	p := parser.New(l, f.Name, f.Abs)
	if p.HasErrors() {
		return nil, p.Errors()[0], nil
	}

	prog := p.ParseProgram()
	prog.AbsPath = f.Abs
	prog.Name = f.Name

	if p.HasErrors() {
		return nil, p.Errors()[0], nil
	}

	return prog, nil, nil
}
