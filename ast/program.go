package ast

import (
	"bytes"
	"fmt"
)

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}

	return ""
}

func (p *Program) String() string {
	var result bytes.Buffer

	for _, stmt := range p.Statements {
		result.WriteString(stmt.String())
	}

	return result.String()
}

func (p *Program) Inserts() map[string]*InsertStatement {
	inserts := make(map[string]*InsertStatement)

	for _, stmt := range p.Statements {
		if insertStmt, ok := stmt.(*InsertStatement); ok {
			inserts[insertStmt.Name.Value] = insertStmt
		}
	}

	return inserts
}

func (p *Program) ApplyInserts(inserts map[string]*InsertStatement) error {
	for _, stmt := range p.Statements {
		if stmt.TokenLiteral() != "reserve" {
			continue
		}

		reserve, ok := stmt.(*ReserveStatement)

		if !ok {
			return fmt.Errorf("expected *ReserveStatement, got %T", stmt)
		}

		insert, hasInsert := inserts[reserve.Name.Value]

		if !hasInsert {
			return fmt.Errorf("The insert statement named '%s' is not defined", reserve.Name.Value)
		}

		reserve.Insert = insert
	}

	return nil
}

// Remove all the statements except the first one, which is the layout statement
// When we use layout, we do not print the file itself
func (p *Program) ApplyLayout(layoutProg *Program) error {
	p.Statements = []Statement{p.Statements[0]}
	p.Statements[0].(*LayoutStatement).Program = layoutProg

	return nil
}
