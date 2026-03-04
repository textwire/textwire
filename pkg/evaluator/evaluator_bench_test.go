package evaluator

import (
	"testing"

	"github.com/textwire/textwire/v3/pkg/lexer"
	"github.com/textwire/textwire/v3/pkg/object"
	"github.com/textwire/textwire/v3/pkg/parser"
)

var inp = `<div>
	@for(i = 0; i < 100000; i++)
		{{ first = 0 }}
		{{ second = 3223230 }}
		{{ third = 0 }}
		{{ forth = 3 }}

		{{ i }}
		{{ first = first + 1000 }}
		{{ second = second / 1 }}
		{{ third = third - 1000 }}
		{{ forth = forth * 2 }}

		@if(first && second && third || forth)
			<h2>HERE</h2>
		@end
	@end
</div>`

func BenchmarkEvaluator(b *testing.B) {
	l := lexer.New(inp)

	p := parser.New(l, nil)
	prog := p.ParseProgram()

	if p.HasErrors() {
		b.Fatal(p.Errors()[0])
	}

	if prog == nil {
		b.Fatal("prog is nil")
	}

	e := New(nil, nil)
	ctx := NewContext(object.NewScope(), prog.AbsPath)

	b.ResetTimer()
	for b.Loop() {
		evaluated := e.Eval(prog, ctx)
		if err, ok := evaluated.(*object.Error); ok {
			b.Fatalf("evaluated is error: %v", err)
		}
	}
}
