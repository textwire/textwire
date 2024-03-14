package ast

import (
	"bytes"

	"github.com/textwire/textwire/fail"
	"github.com/textwire/textwire/token"
)

type Program struct {
	Token      token.Token // The 'program' token
	IsLayout   bool
	UseStmt    *UseStmt
	Statements []Statement
	Components []*ComponentStmt
	Reserves   map[string]*ReserveStmt
	Inserts    map[string]*InsertStmt
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}

	return ""
}

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

func (p *Program) Line() uint {
	return p.Token.Line
}

func (p *Program) ApplyInserts(inserts map[string]*InsertStmt, absPath string) *fail.Error {
	for _, reserve := range p.Reserves {
		insert, hasInsert := inserts[reserve.Name.Value]

		if !hasInsert {
			return fail.New(reserve.Line(), absPath, "parser",
				fail.ErrInsertStmtNotDefined, reserve.Name.Value)
		}

		reserve.Insert = insert
	}

	return nil
}

func (p *Program) ApplyLayout(prog *Program) {
	p.UseStmt.Program = prog
	p.Statements = []Statement{p.UseStmt}
}

func (p *Program) ApplyComponent(name string, prog *Program) {
	for _, comp := range p.Components {
		if comp.Name.Value == name {
			comp.Block = prog
		}
	}
}

func (p *Program) HasReserveStmt() bool {
	return len(p.Reserves) > 0
}

func (p *Program) HasUseStmt() bool {
	return p.UseStmt != nil
}
