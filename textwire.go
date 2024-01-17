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
func ParseFile(filePath string) (*View, error) {
	fullPath, err := getFullPath(filePath)

	if err != nil {
		return nil, err
	}

	content, err := os.ReadFile(fullPath)

	if err != nil {
		return nil, err
	}

	program, err := ParseStr(string(content))

	if err != nil {
		return nil, err
	}

	return &View{program: program}, nil
}

func NewConfig(c *Config) {
	c.TemplateDir = strings.Trim(c.TemplateDir, "/")
	config = c
}
