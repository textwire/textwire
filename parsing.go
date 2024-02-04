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
		return nil, fail.FromError(err, 0, absPath, "template").ToSlice()
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

		if hasLayout, useStmt := prog.HasUseStmt(); hasLayout {
			errs = applyLayoutToProgram(useStmt.Name.Value, prog)
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
	layoutAbsAPath, err := getFullPath(layoutName, true)

	if err != nil {
		return fail.FromError(err, 0, "", "template").ToSlice()
	}

	layoutProg, errs := parseProgram(layoutAbsAPath)
	layoutProg.IsLayout = true

	if len(errs) != 0 {
		return errs
	}

	layoutErr := layoutProg.ApplyInserts(prog.Inserts(), layoutAbsAPath)

	if layoutErr != nil {
		return layoutErr.ToSlice()
	}

	prog.ApplyLayout(layoutProg)

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
