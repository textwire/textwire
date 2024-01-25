package textwire

import (
	"github.com/textwire/textwire/ast"
	"github.com/textwire/textwire/lexer"
	"github.com/textwire/textwire/parser"
)

var config = &Config{
	TemplateDir: "templates",
	TemplateExt: ".textwire.html",
}

// Config is the main configuration for Textwire
type Config struct {
	// TemplateDir is the directory where the Textwire
	// templates are located
	TemplateDir string

	// TemplateExt is the extension of the Textwire
	// template files
	TemplateExt string
}

func ParseStr(text string) (*ast.Program, error) {
	lex := lexer.New(text)
	pars := parser.New(lex)

	prog := pars.ParseProgram()

	if len(pars.Errors()) != 0 {
		return nil, pars.Errors()[0]
	}

	return prog, nil
}

func New(c *Config) (*Template, error) {
	applyConfig(c)

	paths, err := findTextwireFiles()

	if err != nil {
		return nil, err
	}

	programs, err := parsePrograms(paths)

	return &Template{programs: programs}, err
}
