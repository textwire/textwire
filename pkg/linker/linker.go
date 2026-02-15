package linker

import (
	"sync"

	"github.com/textwire/textwire/v3/pkg/ast"
	"github.com/textwire/textwire/v3/pkg/fail"
)

// NodeLinker handles connecting AST nodes between each other to prepare AST
// for evaluator. It will connect @insert to @reserve, @use to layout file,
// @component to its corresponding component file, etc.
type NodeLinker struct {
	Programs []*ast.Program
	mu       sync.RWMutex
}

func New(progs []*ast.Program) *NodeLinker {
	if progs == nil {
		progs = make([]*ast.Program, 0, 4)
	}
	return &NodeLinker{Programs: progs}
}

// LinkNodes links components and layouts to those programs that use them.
// For example, we need to add component program to @component('book'), where
// CompProg is the parsed program AST of the `book.tw` component.
func (nl *NodeLinker) LinkNodes() *fail.Error {
	for _, prog := range nl.Programs {
		if err := nl.handleLayoutLinking(prog); err != nil {
			return err
		}
	}

	for _, prog := range nl.Programs {
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
	layoutProg := ast.FindProg(layoutName, nl.Programs)
	if layoutProg == nil {
		return fail.New(prog.Line(), prog.AbsPath, "API", fail.ErrUseStmtMissingLayout, layoutName)
	}

	layoutProg.IsLayout = true
	if err := ast.CheckUnusedInserts(layoutProg, prog.Inserts); err != nil {
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
		compProg := ast.FindProg(compName, nl.Programs)
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

// Lock acquires the write lock for updating Programs.
func (nl *NodeLinker) Lock() {
	nl.mu.Lock()
}

// Unlock releases the write lock.
func (nl *NodeLinker) Unlock() {
	nl.mu.Unlock()
}

// RLock acquires the read lock for accessing Programs.
func (nl *NodeLinker) RLock() {
	nl.mu.RLock()
}

// RUnlock releases the read lock.
func (nl *NodeLinker) RUnlock() {
	nl.mu.RUnlock()
}
