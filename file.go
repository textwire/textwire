package textwire

import (
	"github.com/textwire/textwire/v3/ast"
	"github.com/textwire/textwire/v3/fail"
	"github.com/textwire/textwire/v3/lexer"
	"github.com/textwire/textwire/v3/parser"
)

// textwireFile holds information about individual Textwire file, including
// relative and absolute file paths.
type textwireFile struct {
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

func NewTextwireFile(rel, abs string) *textwireFile {
	rel = addTwExtension(rel)
	abs = addTwExtension(abs)

	return &textwireFile{
		Rel: trimRelPath(rel),
		Abs: abs,
	}
}

// parseProgram parses twFile.Program and returns errors.
func (tf *textwireFile) parseProgram() (*fail.Error, error) {
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
