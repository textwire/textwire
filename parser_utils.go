package textwire

import (
	"errors"
	"os"

	"github.com/textwire/textwire/v3/ast"
	"github.com/textwire/textwire/v3/fail"
	"github.com/textwire/textwire/v3/lexer"
	"github.com/textwire/textwire/v3/parser"
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

// parseProgram returns program, Textwire error and native error
func parseProgram(absPath, relPath string) (*ast.Program, *fail.Error, error) {
	content, err := fileContent(relPath)
	if err != nil {
		return nil, nil, err
	}

	lex := lexer.New(content)
	pars := parser.New(lex, absPath)
	prog := pars.ParseProgram()

	if len(pars.Errors()) != 0 {
		return nil, pars.Errors()[0], nil
	}

	return prog, nil, nil
}

func parsePrograms(paths map[string]string) (map[string]*ast.Program, *fail.Error) {
	var result = map[string]*ast.Program{}

	for progRelPath, progAbsPath := range paths {
		prog, failure, parseErr := parseProgram(progAbsPath, progRelPath)
		if parseErr != nil {
			return nil, fail.FromError(parseErr, 0, progAbsPath, "template")
		}

		if failure != nil {
			return nil, failure
		}

		if err := applyLayoutToProgram(prog); err != nil {
			return nil, err
		}

		if err := applyComponentToProgram(prog, progAbsPath); err != nil {
			return nil, err
		}

		if !prog.HasReserveStmt() {
			result[progRelPath] = prog
		}
	}

	return result, nil
}

func applyLayoutToProgram(prog *ast.Program) *fail.Error {
	if !prog.HasUseStmt() {
		return nil
	}

	relPath := nameToRelPath(prog.UseStmt.Name.Value)
	absPath, err := getFullPath(relPath)
	if err != nil {
		return fail.FromError(err, prog.UseStmt.Line(), absPath, "template")
	}

	layoutProg, failure, parseErr := parseProgram(absPath, relPath)
	if parseErr != nil {
		return fail.FromError(parseErr, prog.UseStmt.Line(), absPath, "template")
	}

	if failure != nil {
		return failure
	}

	layoutProg.IsLayout = true

	layoutErr := layoutProg.ApplyInserts(prog.Inserts, absPath)
	if layoutErr != nil {
		return layoutErr
	}

	prog.ApplyLayout(layoutProg)

	return nil
}

func applyComponentToProgram(prog *ast.Program, progAbsPath string) *fail.Error {
	for _, comp := range prog.Components {
		relPath := nameToRelPath(comp.Name.Value)

		absPath, err := getFullPath(relPath)
		if err != nil {
			return fail.FromError(err, 0, "", "template")
		}

		compProg, failure, parseErr := parseProgram(absPath, relPath)
		if parseErr != nil {
			if errors.Is(parseErr, os.ErrNotExist) {
				return fail.New(
					comp.Line(),
					progAbsPath,
					"template",
					fail.ErrUndefinedComponent,
					comp.Name.Value,
				)
			}

			return fail.FromError(parseErr, comp.Line(), absPath, "template")
		}

		if failure != nil {
			return failure
		}

		if err := prog.ApplyComponent(comp.Name.Value, compProg, progAbsPath); err != nil {
			return err
		}
	}

	return nil
}
