package textwire

import (
	"strings"

	"github.com/textwire/textwire/ast"
	"github.com/textwire/textwire/fail"
	"github.com/textwire/textwire/lexer"
	"github.com/textwire/textwire/parser"
)

func parseStr(text string) (*ast.Program, []*fail.Error) {
	lex := lexer.New(text)
	pars := parser.New(lex, "")

	prog := pars.ParseProgram()

	if pars.HasErrors() {
		return nil, pars.Errors()
	}

	return prog, nil
}

func parseProgram(absPath string) (*ast.Program, []*fail.Error) {
	content, err := fileContent(absPath)

	if err != nil {
		return nil, []*fail.Error{
			fail.New(0, absPath, "template", err.Error()),
		}
	}

	lex := lexer.New(content)
	pars := parser.New(lex, absPath)

	prog := pars.ParseProgram()

	if len(pars.Errors()) != 0 {
		return nil, pars.Errors()
	}

	return prog, nil
}

func parsePrograms(paths map[string]string) (map[string]*ast.Program, []*fail.Error) {
	var result = map[string]*ast.Program{}

	for name, absPath := range paths {
		prog, errs := parseProgram(absPath)

		if len(errs) != 0 {
			return nil, errs
		}

		if hasLayout, layout := prog.HasUseStmt(); hasLayout {
			errs = applyLayoutToProgram(layout.Name.Value, prog)
		}

		if len(errs) != 0 {
			return nil, errs
		}

		if !prog.HasReserveStmt() {
			result[name] = prog
		}
	}

	return result, nil
}

func applyLayoutToProgram(layoutName string, prog *ast.Program) []*fail.Error {
	layoutAbsAPath, err := getFullPath(layoutName)

	if err != nil {
		return []*fail.Error{
			fail.New(0, layoutAbsAPath, "template", err.Error()),
		}
	}

	layoutProg, errs := parseProgram(layoutAbsAPath)

	if len(errs) != 0 {
		return errs
	}

	err = layoutProg.ApplyInserts(prog.Inserts())

	if err != nil {
		return []*fail.Error{
			fail.New(0, layoutAbsAPath, "template", err.Error()),
		}
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
