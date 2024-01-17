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
	pars := parser.New(lex)

	program := pars.ParseProgram()

	if len(pars.Errors()) != 0 {
		return nil, pars.Errors()[0]
	}

	return program, nil
}

// ParseFile parses a Textwire file and caches the result
func ParseTemplate(filePath string) (*Template, error) {
	program, err := parseProgram(filePath)

	if err != nil {
		return nil, err
	}

	layout, isLayout := program.Statements[0].(*ast.LayoutStatement)

	if isLayout {
		layoutProgram, err := parseProgram(layout.Path.Value)

		if err != nil {
			return nil, err
		}

		program.Statements[0].(*ast.LayoutStatement).Program = layoutProgram
	}

	return &Template{program: program}, nil
}

func NewConfig(c *Config) {
	c.TemplateDir = strings.Trim(c.TemplateDir, "/")
	config = c
}

func parseProgram(filePath string) (*ast.Program, error) {
	content, err := fileContent(filePath)

	if err != nil {
		return nil, err
	}

	lex := lexer.New(content)
	pars := parser.New(lex)

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
