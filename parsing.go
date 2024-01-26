package textwire

import (
	"strings"

	"github.com/textwire/textwire/ast"
	"github.com/textwire/textwire/lexer"
	"github.com/textwire/textwire/parser"
)

func parseProgram(absPath string) (*ast.Program, error) {
	content, err := fileContent(absPath)

	if err != nil {
		return nil, err
	}

	lex := lexer.New(content)
	pars := parser.New(lex)

	prog := pars.ParseProgram()

	if len(pars.Errors()) != 0 {
		return nil, pars.CombinedErrors()
	}

	return prog, nil
}

func parsePrograms(paths map[string]string) (map[string]*ast.Program, error) {
	var result = map[string]*ast.Program{}

	for name, absPath := range paths {
		prog, err := parseProgram(absPath)

		if err != nil {
			return nil, err
		}

		if hasLayout, layout := prog.HasLayout(); hasLayout {
			err = applyLayoutToProgram(layout.Name.Value, prog)
		}

		if err != nil {
			return nil, err
		}

		if !prog.IsLayout() {
			result[name] = prog
		}
	}

	return result, nil
}

func applyLayoutToProgram(layoutName string, prog *ast.Program) error {
	layoutAbsAPath, err := getFullPath(layoutName)

	if err != nil {
		return err
	}

	layoutProg, err := parseProgram(layoutAbsAPath)

	if err != nil {
		return err
	}

	err = layoutProg.ApplyInserts(prog.Inserts())

	if err != nil {
		return err
	}

	prog.ApplyLayout(layoutProg)

	return nil
}

func applyConfig(c *Config) {
	if c.TemplateDir != "" {
		config.TemplateDir = strings.Trim(c.TemplateDir, "/")
	}

	if c.TemplateExt != "" {
		config.TemplateExt = c.TemplateExt
	}
}
