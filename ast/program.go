package ast

import (
	"bytes"

	fail "github.com/textwire/textwire/v3/fail"
	token "github.com/textwire/textwire/v3/token"
)

type Program struct {
	BaseNode
	IsLayout   bool
	Filepath   string
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
	var out bytes.Buffer

	for _, stmt := range p.Statements {
		if _, ok := stmt.(*ExpressionStmt); ok {
			out.WriteString("{{ " + stmt.String() + " }}")
			continue
		}

		out.WriteString(stmt.String())
	}

	return out.String()
}

func (p *Program) Stmts() []Statement {
	res := make([]Statement, 0)

	if p.Statements == nil {
		return []Statement{}
	}

	for _, stmt := range p.Statements {
		if stmt == nil {
			continue
		}

		if s, ok := stmt.(NodeWithStatements); ok {
			res = append(res, s.(Statement))
			res = append(res, s.Stmts()...)
		}
	}

	return res
}

func (p *Program) ApplyInserts(inserts map[string]*InsertStmt, absPath string) *fail.Error {
	if err := p.checkUndefinedInsert(inserts); err != nil {
		return err
	}

	for _, reserve := range p.Reserves {
		insert, hasInsert := inserts[reserve.Name.Value]

		if hasInsert {
			reserve.Insert = insert
		}
	}

	return nil
}

func (p *Program) ApplyLayout(layoutProg *Program) {
	p.UseStmt.Layout = layoutProg
	p.Statements = []Statement{p.UseStmt}
}

func (p *Program) ApplyComponent(name string, prog *Program, progFilePath string) *fail.Error {
	for _, comp := range p.Components {
		if comp.Name.Value != name {
			continue
		}

		duplicateName, times := findDuplicateSlot(comp.Slots)

		if times > 0 {
			if name == "" {
				return fail.New(prog.Line(), progFilePath, "parser",
					fail.ErrDuplicateDefaultSlotUsage, times, name)
			}

			return fail.New(prog.Line(), progFilePath, "parser",
				fail.ErrDuplicateSlotUsage, duplicateName, times, name)
		}

		for _, slot := range comp.Slots {
			idx := findSlotStmtIndex(prog.Statements, slot.Name.Value)

			if idx == -1 {
				if slot.Name.Value == "" {
					return fail.New(prog.Line(), progFilePath, "parser",
						fail.ErrDefaultSlotNotDefined, name)
				}

				return fail.New(prog.Line(), progFilePath, "parser",
					fail.ErrSlotNotDefined, slot.Name.Value, name)
			}

			prog.Statements[idx].(*SlotStmt).Body = slot.Body
		}

		comp.Block = prog
	}

	return nil
}

func (p *Program) HasReserveStmt() bool {
	return len(p.Reserves) > 0
}

func (p *Program) HasUseStmt() bool {
	return p.UseStmt != nil
}

func (p *Program) checkUndefinedInsert(inserts map[string]*InsertStmt) *fail.Error {
	for name := range inserts {
		if _, ok := p.Reserves[name]; ok {
			continue
		}

		line := inserts[name].Line()
		path := inserts[name].FilePath
		name := inserts[name].Name.Value

		return fail.New(line, path, "parser", fail.ErrUndefinedInsert, name)
	}

	return nil
}
