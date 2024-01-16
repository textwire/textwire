package textwire

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/textwire/textwire/evaluator"
	"github.com/textwire/textwire/lexer"
	"github.com/textwire/textwire/object"
	"github.com/textwire/textwire/parser"
)

var config *Config

type Config struct {
	TemplateDir string
}

// ParseStr parses a Textwire string and returns the result as a string
func ParseStr(text string, vars map[string]interface{}) (string, error) {
	lex := lexer.New(text)
	pars := parser.New(lex)

	program := pars.ParseProgram()

	if len(pars.Errors()) != 0 {
		return "", pars.Errors()[0]
	}

	env, err := object.EnvFromMap(vars)

	if err != nil {
		return "", err
	}

	evaluated := evaluator.Eval(program, env)

	if evaluated.Type() == object.ERROR_OBJ {
		return "", errors.New(evaluated.String())
	}

	return evaluated.String(), nil
}

// ParseFile parses a Textwire file and returns the result as a string
func ParseFile(filePath string, vars map[string]interface{}) (string, error) {
	content, err := os.ReadFile(filePath)

	if err != nil {
		return "", err
	}

	strContent := string(content)

	return ParseStr(strContent, vars)
}

func View(w http.ResponseWriter, fileName string, vars map[string]interface{}) error {
	fullPath, err := getFullPath(fileName)

	if err != nil {
		return err
	}

	content, err := ParseFile(fullPath, vars)

	if err != nil {
		return err
	}

	fmt.Fprint(w, content)

	return nil
}

func SetConfig(c *Config) {
	if c.TemplateDir[len(c.TemplateDir)-1:] == "/" {
		c.TemplateDir = c.TemplateDir[:len(c.TemplateDir)-1]
	}

	config = c
}

func getFullPath(fileName string) (string, error) {
	path := fmt.Sprintf("%s/%s.textwire.html", config.TemplateDir, fileName)
	absPath, err := filepath.Abs(path)

	if err != nil {
		return "", err
	}

	return absPath, nil
}
