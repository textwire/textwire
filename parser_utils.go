package textwire

import (
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

// parsePrograms parses each Textwire file into AST nodes.
func parsePrograms(twFiles []*textwireFile) *fail.Error {
	for _, twFile := range twFiles {
		failure, parseErr := twFile.parseProgram()
		if parseErr != nil {
			return fail.FromError(parseErr, 0, twFile.Abs, "template")
		}

		if failure != nil {
			return failure
		}
	}

	return nil
}

// addAttachments adds components and layouts to those files that use them.
// For example, we need to add Attachment to @component('name'), where
// attachment is the parsed AST of the name.tw component.
func addAttachments(twFiles []*textwireFile) *fail.Error {
	for _, twFile := range twFiles {
		if err := addAttachToUseStmt(twFile, twFiles); err != nil {
			return err
		}

		if err := addAttachToCompStmts(twFile, twFiles); err != nil {
			return err
		}
	}

	return nil
}

func addAttachToUseStmt(twFile *textwireFile, twFiles []*textwireFile) *fail.Error {
	if !twFile.Prog.HasUseStmt() {
		return nil
	}

	layoutName := twFile.Prog.UseStmt.Name.Value
	layoutTwFile := findTwFile(layoutName, twFiles)
	if layoutTwFile == nil {
		return fail.New(twFile.Prog.Line(), twFile.Abs, "API", fail.ErrProgramNotFound, layoutName)
	}

	layoutTwFile.Prog.IsLayout = true

	err := layoutTwFile.Prog.AddInsertsAttachments(twFile.Prog.Inserts)
	if err != nil {
		return err
	}

	twFile.Prog.AddLayoutAttachment(layoutTwFile.Prog)

	return nil
}

func addAttachToCompStmts(twFile *textwireFile, twFiles []*textwireFile) *fail.Error {
	if len(twFile.Prog.Components) == 0 {
		return nil
	}

	for _, comp := range twFile.Prog.Components {
		compName := comp.Name.Value
		compTwFile := findTwFile(compName, twFiles)
		if compTwFile == nil {
			return fail.New(
				twFile.Prog.Line(),
				twFile.Abs,
				"API",
				fail.ErrUndefinedComponent,
				compName,
			)
		}

		err := twFile.Prog.AddCompAttachment(compName, compTwFile.Prog, twFile.Abs)
		if err != nil {
			return err
		}
	}

	return nil
}
