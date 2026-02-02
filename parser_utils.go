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

func parsePrograms(twFiles []*textwireFile) *fail.Error {
	for _, twFile := range twFiles {
		failure, parseErr := twFile.parseProgram()
		if parseErr != nil {
			return fail.FromError(parseErr, 0, twFile.Abs, "template")
		}

		if failure != nil {
			return failure
		}

		if err := applyLayoutToProgram(twFile); err != nil {
			return err
		}

		if err := applyComponentToProgram(twFile); err != nil {
			return err
		}
	}

	return nil
}

func applyLayoutToProgram(progTwFile *textwireFile) *fail.Error {
	if !progTwFile.Prog.HasUseStmt() {
		return nil
	}

	layoutRelPath := nameToRelPath(progTwFile.Prog.UseStmt.Name.Value)
	layoutAbsPath, err := getFullPath(layoutRelPath)
	if err != nil {
		return fail.FromError(err, progTwFile.Prog.UseStmt.Line(), layoutAbsPath, "template")
	}

	layoutTwFile := NewTextwireFile(layoutRelPath, layoutAbsPath)

	failure, parseErr := layoutTwFile.parseProgram()
	if parseErr != nil {
		return fail.FromError(parseErr, progTwFile.Prog.UseStmt.Line(), layoutAbsPath, "template")
	}

	if failure != nil {
		return failure
	}

	layoutTwFile.Prog.IsLayout = true

	layoutErr := layoutTwFile.Prog.ApplyInserts(progTwFile.Prog.Inserts, layoutAbsPath)
	if layoutErr != nil {
		return layoutErr
	}

	progTwFile.Prog.ApplyLayout(layoutTwFile.Prog)

	return nil
}

func applyComponentToProgram(progTwFile *textwireFile) *fail.Error {
	for _, comp := range progTwFile.Prog.Components {
		relPath := nameToRelPath(comp.Name.Value)

		absPath, err := getFullPath(relPath)
		if err != nil {
			return fail.FromError(err, 0, "", "template")
		}

		compTwFile := NewTextwireFile(relPath, absPath)

		failure, parseErr := compTwFile.parseProgram()
		if parseErr != nil {
			if errors.Is(parseErr, os.ErrNotExist) {
				return fail.New(
					comp.Line(),
					progTwFile.Abs,
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

		if err := progTwFile.Prog.ApplyComponent(comp.Name.Value, compTwFile.Prog, progTwFile.Abs); err != nil {
			return err
		}
	}

	return nil
}
