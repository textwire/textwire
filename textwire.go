package textwire

import (
	"os"
	"strings"

	"github.com/textwire/textwire/ast"
	"github.com/textwire/textwire/lexer"
	"github.com/textwire/textwire/parser"
)

var config *Config

// Config is the main configuration for Textwire
type Config struct {
	// TemplateDir is the directory where the Textwire
	// templates are located
	TemplateDir string
}

// ParseStr parses a Textwire string and returns the result as a string
func ParseStr(text string) (*ast.Program, error) {
	lex := lexer.New(text)
	pars := parser.New(lex, nil)

	program := pars.ParseProgram()

	if len(pars.Errors()) != 0 {
		return nil, pars.Errors()[0]
	}

	return program, nil
}

// ParseFile parses a Textwire file and caches the result
func ParseTemplate(filePath string) (*Template, error) {
	prog, err := parseProgram(filePath, nil)

	if err != nil {
		return nil, err
	}

	layout, isLayout := prog.Statements[0].(*ast.LayoutStatement)

	if !isLayout {
		return &Template{program: prog}, nil
	}

	// Parse the layout program
	layoutProg, err := parseProgram(layout.Path.Value, prog.Inserts())

	if err != nil {
		return nil, err
	}

	prog.Statements[0].(*ast.LayoutStatement).Program = layoutProg

	return &Template{program: prog}, nil
}

func NewConfig(c *Config) {
	c.TemplateDir = strings.Trim(c.TemplateDir, "/")
	config = c
}

func parseProgram(filePath string, inserts map[string]*ast.InsertStatement) (*ast.Program, error) {
	content, err := fileContent(filePath)

	if err != nil {
		return nil, err
	}

	lex := lexer.New(content)
	pars := parser.New(lex, inserts)

	program := pars.ParseProgram()

	if len(pars.Errors()) != 0 {
		return nil, pars.Errors()[0]
	}

	return program, nil
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
