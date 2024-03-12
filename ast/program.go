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
	for _, stmt := range p.Statements {
		if stmt.TokenLiteral() != token.String(token.RESERVE) {
			continue
		}

		reserve, ok := stmt.(*ReserveStmt)

		if !ok {
			return fail.New(stmt.Line(), absPath, "parser",
				fail.ErrExceptedReserveStmt, stmt)
		}

		insert, hasInsert := inserts[reserve.Name.Value]

		if !hasInsert {
			return fail.New(stmt.Line(), absPath, "parser",
				fail.ErrInsertStmtNotDefined, reserve.Name.Value)
		}

		reserve.Insert = insert
	}

	return nil
}

// Remove all the statements except the first one, which is the layout statement
// When we use layout, we do not print the file itself
func (p *Program) ApplyLayout(prog *Program) {
	p.Statements = []Statement{p.Statements[0]}
	p.Statements[0].(*UseStmt).Program = prog
}

func (p *Program) ApplyComponent(name string, prog *Program) {
	for _, comp := range p.Components {
		if comp.Name.Value == name {
			comp.Block = prog
		}
	}
}

func (p *Program) HasReserveStmt() bool {
	for _, stmt := range p.Statements {
		if stmt.TokenLiteral() == token.String(token.RESERVE) {
			return true
		}
	}

	return false
}

func (p *Program) HasUseStmt() bool {
	return p.UseStmt != nil
}

func (p *Program) HasComponentStmt() bool {
	return len(p.Components) > 0
}
