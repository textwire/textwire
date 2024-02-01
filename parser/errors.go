package parser

import (
	"github.com/textwire/textwire/fail"
)

func (p *Parser) Errors() []*fail.Error {
	return p.errors
}

func (p *Parser) HasErrors() bool {
	return len(p.errors) > 0
}
