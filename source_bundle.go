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

// SourceBundle is the main struct to handle parsing and evaluation of
// Textwire code.
type SourceBundle struct {
	files    []*file
	programs []*ast.Program
}

func NewSourceBundle() *SourceBundle {
	const approximateCap = 4

	return &SourceBundle{
		files:    make([]*file, 0, approximateCap),
		programs: make([]*ast.Program, 0, approximateCap),
	}
}

// ParseFiles parses each Textwire file into AST nodes and saves them.
func (sb *SourceBundle) ParseFiles() *fail.Error {
	for _, f := range sb.files {
		prog, failure, parseErr := sb.parseFile(f)
		if parseErr != nil {
			return fail.FromError(parseErr, 0, f.Abs, "template")
		}

		if failure != nil {
			return failure
		}

		sb.programs = append(sb.programs, prog)
	}

	return nil
}

// AddAttachments adds components and layouts to those programs that use them.
// For example, we need to add Attachment to @component('book'), where
// attachment is the parsed program AST of the book.tw component.
func (sb *SourceBundle) AddAttachments() *fail.Error {
	for _, prog := range sb.programs {
		if err := sb.addAttachToUse(prog); err != nil {
			return err
		}

		if err := sb.addAttachToComp(prog); err != nil {
			return err
		}
	}

	return nil
}

func (sb *SourceBundle) FindProg(name string) *ast.Program {
	for i := range sb.programs {
		if sb.programs[i].Name == name {
			return sb.programs[i]
		}
	}

	return nil
}

// FindFiles recursively finds all Textwire files in the templates directory,
// and creates a *file wrapper for each of these files.
func (sb *SourceBundle) FindFiles() error {
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

// addAttachToUse add attachments to use statement
func (sb *SourceBundle) addAttachToUse(prog *ast.Program) *fail.Error {
	if !prog.HasUseStmt() {
		return nil
	}

	layoutName := prog.UseStmt.Name.Value
	layoutProg := sb.FindProg(layoutName)
	if layoutProg == nil {
		return fail.New(prog.Line(), prog.AbsPath, "API", fail.ErrUseStmtMissingLayout, layoutName)
	}

	layoutProg.IsLayout = true
	err := layoutProg.AddInsertsAttachments(prog.Inserts)
	if err != nil {
		return err
	}

	prog.AddLayoutAttachment(layoutProg)

	return nil
}

// addAttachToComp adds attachments to components
func (sb *SourceBundle) addAttachToComp(prog *ast.Program) *fail.Error {
	if len(prog.Components) == 0 {
		return nil
	}

	for _, comp := range prog.Components {
		compName := comp.Name.Value
		compProg := sb.FindProg(compName)
		if compProg == nil {
			return fail.New(prog.Line(), prog.AbsPath, "API", fail.ErrUndefinedComponent, compName)
		}

		err := prog.AddCompAttachment(compName, compProg, prog.AbsPath)
		if err != nil {
			return err
		}
	}

	return nil
}

// parseFile parses given file into a ast.Program and returns it.
func (sb *SourceBundle) parseFile(f *file) (*ast.Program, *fail.Error, error) {
	content, err := f.Content()
	if err != nil {
		return nil, nil, err
	}

	l := lexer.New(content)
	p := parser.New(l, f.Abs)
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
