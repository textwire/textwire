package ast

import (
	"bytes"
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
