package ast

import (
	"bytes"

	"github.com/textwire/textwire/fail"
	"github.com/textwire/textwire/token"
)

type Program struct {
	Token      token.Token // The 'program' token
	Statements []Statement
	IsLayout   bool
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

func (p *Program) Inserts() map[string]*InsertStmt {
	inserts := make(map[string]*InsertStmt)

	for _, stmt := range p.Statements {
		if insertStmt, ok := stmt.(*InsertStmt); ok {
			inserts[insertStmt.Name.Value] = insertStmt
		}
	}

	return inserts
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
func (p *Program) ApplyLayout(layoutProg *Program) error {
	p.Statements = []Statement{p.Statements[0]}
	p.Statements[0].(*UseStatement).Program = layoutProg

	return nil
}

func (p *Program) HasReserveStmt() bool {
	for _, stmt := range p.Statements {
		if stmt.TokenLiteral() == token.String(token.RESERVE) {
			return true
		}
	}

	return false
}

func (p *Program) HasUseStmt() (bool, *UseStatement) {
	for _, stmt := range p.Statements {
		if stmt.TokenLiteral() != token.String(token.USE) {
			continue
		}

		return true, stmt.(*UseStatement)
	}

	return false, nil
}
