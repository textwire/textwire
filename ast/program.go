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
		_, isHtmlStmt := stmt.(*HTMLStatement)

		if isHtmlStmt {
			result.WriteString(stmt.String())
		} else {
			result.WriteString("{{ " + stmt.String() + " }}")
		}
	}

	return result.String()
}
