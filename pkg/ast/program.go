package ast

import (
	"strings"

	"github.com/textwire/textwire/v3/pkg/fail"
	"github.com/textwire/textwire/v3/pkg/token"
)

type Program struct {
	BaseNode
	IsLayout   bool
	Name       string
	AbsPath    string
	UseStmt    *UseStmt
	Statements []Statement
	Components []*ComponentStmt
	Reserves   map[string]*ReserveStmt
	Inserts    map[string]*InsertStmt
}

func NewProgram(tok token.Token) *Program {
	return &Program{
		BaseNode: NewBaseNode(tok),
	}
}

func (p *Program) statementNode() {}

func (p *Program) String() string {
	var out strings.Builder
	out.Grow(len(p.Statements))

	for i := range p.Statements {
		out.WriteString(p.Statements[i].String())
	}

	return out.String()
}

func (p *Program) Stmts() []Statement {
	stmts := make([]Statement, 0)
	if p.Statements == nil {
		return []Statement{}
	}

	for _, stmt := range p.Statements {
		if stmt == nil {
			continue
		}

		if s, ok := stmt.(NodeWithStatements); ok {
			stmts = append(stmts, s.(Statement))
			stmts = append(stmts, s.Stmts()...)
		}
	}

	return stmts
}

// LinkLayoutToUse adds Layout AST program to UseStmt for the current template
// and resets statements to UseStmt only. Because we don't need anything else
// inside a template. Make sure inserts are added before this is called
// because they will be removed by this function.
func (p *Program) LinkLayoutToUse(layoutProg *Program) {
	p.UseStmt.LayoutProg = layoutProg
	p.Statements = []Statement{p.UseStmt}
}

func (p *Program) LinkCompProg(compName string, prog *Program, absPath string) *fail.Error {
	for _, comp := range p.Components {
		if comp.Name.Value != compName {
			continue
		}

		duplicate, times := findDuplicateSlot(comp.Slots)
		if times > 0 && duplicate != nil {
			if duplicate.IsDefault() {
				return fail.New(
					prog.Line(),
					absPath,
					"parser",
					fail.ErrDuplicateDefaultSlot,
					times,
					compName,
				)
			}

			return fail.New(
				prog.Line(),
				absPath,
				"parser",
				fail.ErrDuplicateSlot,
				duplicate.Name().Value,
				times,
				compName,
			)
		}

		for _, slot := range comp.Slots {
			name := slot.Name().Value
			idx := findSlotIndex(prog.Statements, name)
			if idx != -1 {
				prog.Statements[idx].(SlotStatement).SetBlock(slot.Block())
				continue
			}

			if slot.IsDefault() {
				return fail.New(
					prog.Line(),
					absPath,
					"parser",
					fail.ErrDefaultSlotNotDefined,
					compName,
				)
			}

			return fail.New(prog.Line(), absPath, "parser", fail.ErrSlotNotDefined, compName, name)
		}

		comp.CompProg = prog
	}

	return nil
}

func (p *Program) HasReserveStmt() bool {
	return len(p.Reserves) > 0
}

func (p *Program) HasUseStmt() bool {
	return p.UseStmt != nil
}
