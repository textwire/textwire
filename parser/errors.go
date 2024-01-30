package parser

import (
	"bytes"
	"fmt"

	"github.com/textwire/textwire/fail"
)

func (p *Parser) Errors() []*fail.Error {
	return p.errors
}

func (p *Parser) CombinedErrors() error {
	var out bytes.Buffer

	for _, err := range p.errors {
		out.WriteString(err.String())
		out.WriteString("\n")
	}

	out.Truncate(out.Len() - 1)

	return fmt.Errorf(out.String())
}

func (p *Parser) HasErrors() bool {
	return len(p.errors) > 0
}
