package textwire

import (
	"os"
	"strings"

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

// ParseStr parses a Textwire string and returns the result as a string
func ParseStr(text string) (*ast.Program, error) {
	lex := lexer.New(text)
	pars := parser.New(lex)

	prog := pars.ParseProgram()

	if len(pars.Errors()) != 0 {
		return nil, pars.Errors()[0]
	}

	return prog, nil
}

// ParseFile parses a Textwire file and caches the result
func ParseTemplate(filePath string) (*Template, error) {
	prog, err := parseProgram(filePath)

	if err != nil {
		return nil, err
	}

	layout, isLayout := prog.Statements[0].(*ast.LayoutStatement)

	if !isLayout {
		return &Template{prog: prog}, nil
	}

	// Parse the layout program
	layoutProg, err := parseProgram(layout.Path.Value)

	if err != nil {
		return nil, err
	}

	err = layoutProg.ApplyInserts(prog.Inserts())

	if err != nil {
		return nil, err
	}

	prog.ApplyLayout(layoutProg)

	return &Template{prog: prog}, nil
}

func NewConfig(c *Config) {
	if c.TemplateDir != "" {
		config.TemplateDir = strings.Trim(c.TemplateDir, "/")
	}

	if c.TemplateExt != "" {
		config.TemplateExt = c.TemplateExt
	}
}

func parseProgram(filePath string) (*ast.Program, error) {
	content, err := fileContent(filePath)

	if err != nil {
		return nil, err
	}

	lex := lexer.New(content)
	pars := parser.New(lex)

	prog := pars.ParseProgram()

	if len(pars.Errors()) != 0 {
		return nil, pars.Errors()[0]
	}

	return prog, nil
}

func fileContent(filePath string) (string, error) {
	fullPath, err := getFullPath(filePath)

	if err != nil {
		return "", err
	}

	content, err := os.ReadFile(fullPath)

	if err != nil {
		return "", err
	}

	return string(content), nil
}
