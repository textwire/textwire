package linker

import (
	"github.com/textwire/textwire/v3/ast"
	"github.com/textwire/textwire/v3/fail"
)

// NodeLinker handles connecting AST nodes between each other to prepare AST
// for evaluator. It will connect @insert to @reserve, @use to layout file,
// @component to its corresponding component file, etc.
type NodeLinker struct {
	programs []*ast.Program
}

func New(progs []*ast.Program) *NodeLinker {
	if progs == nil {
		progs = make([]*ast.Program, 0, 4)
	}
	return &NodeLinker{progs}
}

// Progs returns parsed AST Program nodes
func (nl *NodeLinker) Progs() []*ast.Program {
	return nl.programs
}

// LinkNodes links components and layouts to those programs that use them.
// For example, we need to add component program to @component('book'), where
// CompProg is the parsed program AST of the `book.tw` component.
func (nl *NodeLinker) LinkNodes() *fail.Error {
	for _, prog := range nl.programs {
		if err := nl.handleLayoutLinking(prog); err != nil {
			return err
		}
	}

	for _, prog := range nl.programs {
		if err := nl.handleCompLinking(prog); err != nil {
			return err
		}
	}

	return nil
}

// handleLayoutLinking links layout directives to template directives
func (nl *NodeLinker) handleLayoutLinking(prog *ast.Program) *fail.Error {
	if !prog.HasUseStmt() {
		return nil
	}

	prog.UseStmt.Inserts = prog.Inserts

	layoutName := prog.UseStmt.Name.Value
	layoutProg := ast.FindProg(layoutName, nl.programs)
	if layoutProg == nil {
		return fail.New(prog.Line(), prog.AbsPath, "API", fail.ErrUseStmtMissingLayout, layoutName)
	}

	layoutProg.IsLayout = true
	if err := ast.CheckUndefinedInserts(layoutProg, prog.Inserts); err != nil {
		return err
	}

	prog.LinkLayoutToUse(layoutProg)

	return nil
}

// handleCompLinking links component directives with component files
func (nl *NodeLinker) handleCompLinking(prog *ast.Program) *fail.Error {
	if len(prog.Components) == 0 {
		return nil
	}

	for _, comp := range prog.Components {
		compName := comp.Name.Value
		compProg := ast.FindProg(compName, nl.programs)
		if compProg == nil {
			return fail.New(prog.Line(), prog.AbsPath, "API", fail.ErrUndefinedComponent, compName)
		}

		err := prog.LinkCompProg(compName, compProg, prog.AbsPath)
		if err != nil {
			return err
		}
	}

	return nil
}
