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

func parseProgram(absPath string) (*ast.Program, *fail.Error) {
	content, err := fileContent(absPath)

	if err != nil {
		return nil, fail.FromError(err, 0, absPath, "template")
	}

	lex := lexer.New(content)
	pars := parser.New(lex, absPath)
	prog := pars.ParseProgram()

	if len(pars.Errors()) != 0 {
		return nil, pars.Errors()[0]
	}

	return prog, nil
}

func parsePrograms(paths map[string]string) (map[string]*ast.Program, *fail.Error) {
	var result = map[string]*ast.Program{}

	for name, absPath := range paths {
		prog, err := parseProgram(absPath)

		if err != nil {
			return nil, err
		}

		if prog.HasUseStmt() {
			err = applyLayoutToProgram(prog)

			if err != nil {
				return nil, err
			}
		}

		if prog.HasComponentStmt() {
			if err = applyComponentToProgram(prog); err != nil {
				return nil, err
			}
		}

		if !prog.HasReserveStmt() {
			result[name] = prog
		}
	}

	return result, nil
}

func applyLayoutToProgram(prog *ast.Program) *fail.Error {
	layoutName := prog.UseStmt.Name.Value
	layoutAbsPath, err := getFullPath(layoutName, true)

	if err != nil {
		return fail.FromError(err, 0, "", "template")
	}

	layoutProg, parseErr := parseProgram(layoutAbsPath)
	layoutProg.IsLayout = true

	if parseErr != nil {
		return parseErr
	}

	layoutErr := layoutProg.ApplyInserts(prog.Inserts, layoutAbsPath)

	if layoutErr != nil {
		return layoutErr
	}

	prog.ApplyLayout(layoutProg)

	return nil
}

func applyComponentToProgram(prog *ast.Program) *fail.Error {
	for _, comp := range prog.Components {
		compName := comp.Name.Value
		compAbsPath, err := getFullPath(compName, true)

		if err != nil {
			return fail.FromError(err, 0, "", "template")
		}

		compProg, parseErr := parseProgram(compAbsPath)

		if parseErr != nil {
			return parseErr
		}

		prog.ApplyComponent(compName, compProg)
	}

	return nil
}

func applyConfig(c *Config) {
	config.configApplied = true

	if c.TemplateDir != "" {
		config.TemplateDir = strings.Trim(c.TemplateDir, "/")
	}

	if c.TemplateExt != "" {
		config.TemplateExt = c.TemplateExt
	}
}
