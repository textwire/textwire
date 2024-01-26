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
	var msg bytes.Buffer

	for _, err := range p.errors {
		msg.WriteString(err.String())
		msg.WriteString("\n")
	}

	msg.Truncate(msg.Len() - 1)

	return fmt.Errorf(msg.String())
}

func (p *Parser) HasErrors() bool {
	return len(p.errors) > 0
}
