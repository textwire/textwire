package textwire

import (
	"github.com/textwire/textwire/v3/ast"
	"github.com/textwire/textwire/v3/fail"
	"github.com/textwire/textwire/v3/lexer"
	"github.com/textwire/textwire/v3/parser"
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

	// Prog is parsed AST for this Textwire file.
	Prog *ast.Program
}

func NewFile(name, rel, abs string) *file {
	rel = addTwExtension(rel)
	abs = addTwExtension(abs)

	return &file{
		Name: name,
		Rel:  trimRelPath(rel),
		Abs:  abs,
	}
}

// parseProgram parses file.Program and returns errors.
func (tf *file) parseProgram() (*fail.Error, error) {
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
