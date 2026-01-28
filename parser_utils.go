package textwire

import (
	"errors"
	"os"

	"github.com/textwire/textwire/v2/ast"
	"github.com/textwire/textwire/v2/fail"
	"github.com/textwire/textwire/v2/lexer"
	"github.com/textwire/textwire/v2/parser"
)

func parseStr(text string) (*ast.Program, []*fail.Error) {
	lex := lexer.New(text)
	pars := parser.New(lex, "")

	prog := pars.Parse()

	if pars.HasErrors() {
		return nil, pars.Errors()
	}

	return prog, nil
}

// parseProgram returns program, Textwire error and native error
func parseProgram(absPath string) (*ast.Program, *fail.Error, error) {
	content, err := fileContent(absPath)
	if err != nil {
		return nil, nil, err
	}

	lex := lexer.New(content)
	pars := parser.New(lex, absPath)
	prog := pars.Parse()

	if len(pars.Errors()) != 0 {
		return nil, pars.Errors()[0], nil
	}

	return prog, nil, nil
}

func parsePrograms(paths map[string]string) (map[string]*ast.Program, *fail.Error) {
	var result = map[string]*ast.Program{}

	for name, absPath := range paths {
		prog, failure, parseErr := parseProgram(absPath)
		if parseErr != nil {
			return nil, fail.FromError(parseErr, 0, absPath, "template")
		}

		if failure != nil {
			return nil, failure
		}

		if err := applyLayoutToProgram(prog); err != nil {
			return nil, err
		}

		if err := applyComponentToProgram(prog, absPath); err != nil {
			return nil, err
		}

		if !prog.HasReserveStmt() {
			result[name] = prog
		}
	}

	return result, nil
}

func applyLayoutToProgram(prog *ast.Program) *fail.Error {
	if !prog.HasUseStmt() {
		return nil
	}

	layoutName := prog.UseStmt.Name.Value
	stmt := prog.UseStmt
	layoutAbsPath, err := getFullPath(layoutName, true)
	if err != nil {
		return fail.FromError(err, stmt.Line(), layoutAbsPath, "template")
	}

	layoutProg, failure, parseErr := parseProgram(layoutAbsPath)
	if parseErr != nil {
		return fail.FromError(parseErr, stmt.Line(), layoutAbsPath, "template")
	}

	if failure != nil {
		return failure
	}

	layoutProg.IsLayout = true

	layoutErr := layoutProg.ApplyInserts(prog.Inserts, layoutAbsPath)
	if layoutErr != nil {
		return layoutErr
	}

	prog.ApplyLayout(layoutProg)

	return nil
}

func applyComponentToProgram(prog *ast.Program, progFilePath string) *fail.Error {
	for _, comp := range prog.Components {
		compName := comp.Name.Value
		compAbsPath, err := getFullPath(compName, true)
		if err != nil {
			return fail.FromError(err, 0, "", "template")
		}

		compProg, failure, parseErr := parseProgram(compAbsPath)
		if parseErr != nil {
			if errors.Is(parseErr, os.ErrNotExist) {
				return fail.New(comp.Line(), progFilePath, "template",
					fail.ErrUndefinedComponent, comp.Name.Value)
			}

			return fail.FromError(parseErr, comp.Line(), compAbsPath, "template")
		}

		if failure != nil {
			return failure
		}

		if err := prog.ApplyComponent(compName, compProg, progFilePath); err != nil {
			return err
		}
	}

	return nil
}
