package parser

import (
	"testing"

	"codeberg.org/textwire/textwire/v5/pkg/ast"
	"codeberg.org/textwire/textwire/v5/pkg/lexer"
)

func testEmbeddedString[T ast.Segment](t *testing.T, inp, expect string) {
	t.Helper()

	node, err := parseEmbedded[T](inp)
	if err != nil {
		t.Fatal(err)
	}

	if err := testString(node, expect); err != nil {
		t.Fatal(err)
	}
}

func testDirectiveString[T ast.Chunk](t *testing.T, inp, expect string) {
	t.Helper()

	node, err := parseDirective[T](inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	if err := testString(node, expect); err != nil {
		t.Fatal(err)
	}
}

func testChunkInLoopString[T ast.Chunk](t *testing.T, inp, expect string) {
	t.Helper()

	node, err := parseChunkInLoop[T](inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	if err := testString(node, expect); err != nil {
		t.Fatal(err)
	}
}

func TestIdentExprString(t *testing.T) {
	cases := []struct {
		inp    string
		expect string
	}{
		{inp: "{{ x }}", expect: "x"},
		{inp: "{{ myName }}", expect: "myName"},
		{inp: "{{ user_name2 }}", expect: "user_name2"},
	}

	for _, tc := range cases {
		t.Run(tc.inp, func(t *testing.T) {
			testEmbeddedString[*ast.IdentExpr](t, tc.inp, tc.expect)
		})
	}
}

func TestIntExprString(t *testing.T) {
	cases := []struct {
		inp    string
		expect string
	}{
		{inp: "{{ 0 }}", expect: "0"},
		{inp: "{{ 234 }}", expect: "234"},
		{inp: "{{ 1234567890 }}", expect: "1234567890"},
	}

	for _, tc := range cases {
		t.Run(tc.inp, func(t *testing.T) {
			testEmbeddedString[*ast.IntExpr](t, tc.inp, tc.expect)
		})
	}
}

func TestFloatExprString(t *testing.T) {
	cases := []struct {
		inp    string
		expect string
	}{
		{inp: "{{ 0.0 }}", expect: "0.0"},
		{inp: "{{ 2.34149 }}", expect: "2.34149"},
		{inp: "{{ 10.5 }}", expect: "10.5"},
	}

	for _, tc := range cases {
		t.Run(tc.inp, func(t *testing.T) {
			testEmbeddedString[*ast.FloatExpr](t, tc.inp, tc.expect)
		})
	}
}

func TestNilExprString(t *testing.T) {
	testEmbeddedString[*ast.NilExpr](t, "{{ nil }}", "nil")
}

func TestStrExprString(t *testing.T) {
	cases := []struct {
		inp    string
		expect string
	}{
		{inp: "{{ '' }}", expect: `""`},
		{inp: "{{ 'Skirk' }}", expect: `"Skirk"`},
		{inp: `{{ "Kokomi" }}`, expect: `"Kokomi"`},
		{inp: "{{ 'hello world' }}", expect: `"hello world"`},
	}

	for _, tc := range cases {
		t.Run(tc.inp, func(t *testing.T) {
			testEmbeddedString[*ast.StrExpr](t, tc.inp, tc.expect)
		})
	}
}

func TestBoolExprString(t *testing.T) {
	cases := []struct {
		inp    string
		expect string
	}{
		{inp: "{{ true }}", expect: "true"},
		{inp: "{{ false }}", expect: "false"},
	}

	for _, tc := range cases {
		t.Run(tc.inp, func(t *testing.T) {
			testEmbeddedString[*ast.BoolExpr](t, tc.inp, tc.expect)
		})
	}
}

func TestInfixExprString(t *testing.T) {
	testEmbeddedString[*ast.InfixExpr](t, "{{ 3 / 2 }}", "(3 / 2)")
}

func TestPrefixExprString(t *testing.T) {
	cases := []struct {
		inp    string
		expect string
	}{
		{inp: "{{ -234 }}", expect: "(-234)"},
		{inp: "{{ -obj.test }}", expect: "(-(obj.test))"},
		{inp: "{{ !true }}", expect: "(!true)"},
		{inp: "{{ !false }}", expect: "(!false)"},
	}

	for _, tc := range cases {
		t.Run(tc.inp, func(t *testing.T) {
			testEmbeddedString[*ast.PrefixExpr](t, tc.inp, tc.expect)
		})
	}
}

func TestTernaryExprString(t *testing.T) {
	testEmbeddedString[*ast.TernaryExpr](t, "{{ true ? 1 : 0 }}", "(true ? 1 : 0)")
}

func TestArrExprString(t *testing.T) {
	testEmbeddedString[*ast.ArrExpr](t, "{{ [11, 234] }}", "[11, 234]")
}

func TestObjExprString(t *testing.T) {
	testEmbeddedString[*ast.ObjExpr](
		t,
		`{{ {"father": {"name": "John"}} }}`,
		`{"father": {"name": "John"}}`,
	)
}

func TestDotExprString(t *testing.T) {
	testEmbeddedString[*ast.DotExpr](t, "{{ person.father.name }}", "((person.father).name)")
}

func TestIndexExprString(t *testing.T) {
	testEmbeddedString[*ast.IndexExpr](t, "{{ arr[1 + 2][2] }}", "((arr[(1 + 2)])[2])")
}

func TestCallExprString(t *testing.T) {
	testEmbeddedString[*ast.CallExpr](t, "{{ name.len() }}", "(name.len())")
}

func TestGlobalCallExprString(t *testing.T) {
	testEmbeddedString[*ast.GlobalCallExpr](t, "{{ defined(name) }}", "(defined(name))")
}

func TestIncStmtString(t *testing.T) {
	testEmbeddedString[*ast.IncStmt](t, "{{ i++ }}", "(i++)")
}

func TestDecStmtString(t *testing.T) {
	testEmbeddedString[*ast.DecStmt](t, "{{ i-- }}", "(i--)")
}

func TestAssignStmtString(t *testing.T) {
	cases := []struct {
		inp    string
		expect string
	}{
		{inp: `{{ name = "Escoffier" }}`, expect: `(name = "Escoffier")`},
		{inp: `{{ stars = 5 }}`, expect: `(stars = 5)`},
		{inp: `{{ me.age = 34 }}`, expect: `((me.age) = 34)`},
		{inp: `{{ arr[0] = 1 }}`, expect: `((arr[0]) = 1)`},
		{
			inp:    `{{ arr[234][2][23].name.first = "Anna" }}`,
			expect: `((((((arr[234])[2])[23]).name).first) = "Anna")`,
		},
		{inp: `{{ (obj.one.two) = "test" }}`, expect: `(((obj.one).two) = "test")`},
	}

	for _, tc := range cases {
		t.Run(tc.inp, func(t *testing.T) {
			testEmbeddedString[*ast.AssignStmt](t, tc.inp, tc.expect)
		})
	}
}

func TestOperatorPrecedenceString(t *testing.T) {
	cases := []struct {
		id     uint
		inp    string
		expect string
	}{
		{
			id:     10,
			inp:    "{{ 1 * 2 }}",
			expect: "{{ (1 * 2) }}",
		},
		{
			id:     20,
			inp:    "<h2>{{ -2 + 3 }}</h2>",
			expect: "<h2>{{ ((-2) + 3) }}</h2>",
		},
		{
			id:     30,
			inp:    "{{ a + b + c }}",
			expect: "{{ ((a + b) + c) }}",
		},
		{
			id:     40,
			inp:    "{{ a + b / c }}",
			expect: "{{ (a + (b / c)) }}",
		},
		{
			id:     50,
			inp:    "{{ -2.float() }}",
			expect: "{{ (-(2.float())) }}",
		},
		{
			id:     60,
			inp:    "{{ -5.0.int() }}",
			expect: "{{ (-(5.0.int())) }}",
		},
		{
			id:     70,
			inp:    "{{ -obj.test }}",
			expect: "{{ (-(obj.test)) }}",
		},
		{
			id:     80,
			inp:    "{{ true && true || false }}",
			expect: "{{ ((true && true) || false) }}",
		},
		{
			id:     90,
			inp:    "{{ true ? 1 : 0 }}",
			expect: "{{ (true ? 1 : 0) }}",
		},
		{
			id:     100,
			inp:    "{{ true && false ? 1 : 0 }}",
			expect: "{{ ((true && false) ? 1 : 0) }}",
		},
		{
			id:     110,
			inp:    "{{ true && false || 1 ? 1 : 0 }}",
			expect: "{{ (((true && false) || 1) ? 1 : 0) }}",
		},
		{
			id:     120,
			inp:    "{{ -2.float() && -2.0.int() ? 1 : 0 }}",
			expect: "{{ (((-(2.float())) && (-(2.0.int()))) ? 1 : 0) }}",
		},
		{
			id:     130,
			inp:    "{{ !defined(age) || !defined(name) ? 1 : 0 }}",
			expect: "{{ (((!(defined(age))) || (!(defined(name)))) ? 1 : 0) }}",
		},
		{
			id:     140,
			inp:    "{{ defined(name) }}",
			expect: "{{ (defined(name)) }}",
		},
		{
			id:     150,
			inp:    "{{ long = user.name.len() > 0 }}",
			expect: "{{ (long = (((user.name).len()) > 0)) }}",
		},
		{
			id:     160,
			inp:    "{{ user && user.name == 'serhii' }}",
			expect: `{{ (user && ((user.name) == "serhii")) }}`,
		},
	}

	for _, tc := range cases {
		l := lexer.New(tc.inp)
		p := New(l, nil)
		prog := p.ParseProgram()

		if p.HasErrors() {
			t.Fatalf("Case: %d. Parser error:  %v", tc.id, p.Errors()[0])
		}

		if err := testString(prog, tc.expect); err != nil {
			t.Fatalf("Case: %d. %v", tc.id, err)
		}
	}
}

func TestEmbeddedRawString(t *testing.T) {
	testDirectiveString[*ast.Embedded](t, "{!! 'x' !!}", `{!! "x" !!}`)
}

func TestTextString(t *testing.T) {
	inp := "<div><span>Hello</span></div>"
	testDirectiveString[*ast.Text](t, inp, inp)
}

func TestUseDirString(t *testing.T) {
	inp := `@use("main")`
	testDirectiveString[*ast.UseDir](t, inp, inp)
}

func TestReserveDirString(t *testing.T) {
	t.Run("without fallback", func(t *testing.T) {
		inp := `@reserve("content")`
		testDirectiveString[*ast.ReserveDir](t, inp, inp)
	})

	t.Run("with fallback", func(t *testing.T) {
		testDirectiveString[*ast.ReserveDir](
			t,
			`@reserve("content", 'fallback')`,
			`@reserve("content", "fallback")`,
		)
	})
}

func TestInsertDirString(t *testing.T) {
	t.Run("with block", func(t *testing.T) {
		inp := `@insert("content")<h1>Some content</h1>@end`
		testDirectiveString[*ast.InsertDir](t, inp, inp)
	})

	t.Run("with argument", func(t *testing.T) {
		testDirectiveString[*ast.InsertDir](
			t,
			"@insert('content', 'Some content')",
			`@insert("content", "Some content")`,
		)
	})
}

func TestCompDirString(t *testing.T) {
	t.Run("without @pass", func(t *testing.T) {
		testDirectiveString[*ast.CompDir](
			t,
			"@component('components/book-card', { c: card })@end",
			`@component("components/book-card", {"c": card})@end`,
		)
	})

	t.Run("with default @pass", func(t *testing.T) {
		inp := `@component("components/book-card")<h1>Header</h1>@end`
		testDirectiveString[*ast.CompDir](t, inp, inp)
	})

	t.Run("with 2 @pass directives", func(t *testing.T) {
		inp := `@component("x")@pass("header")<h1>Header</h1>@end@pass("footer")<footer>Footer</footer>@end@end`
		testDirectiveString[*ast.CompDir](t, inp, inp)
	})

	t.Run("with @passif", func(t *testing.T) {
		inp := `@component("user")@passif(false, "name")Test2@end@end`
		testDirectiveString[*ast.CompDir](t, inp, inp)
	})
}

func TestSlotDirString(t *testing.T) {
	t.Run("named slot", func(t *testing.T) {
		testDirectiveString[*ast.SlotDir](t, "@slot('header')", `@slot("header")`)
	})

	t.Run("default slot", func(t *testing.T) {
		testDirectiveString[*ast.SlotDir](t, "@slot", "@slot")
	})
}

func TestDumpDirString(t *testing.T) {
	testDirectiveString[*ast.DumpDir](
		t,
		`@dump("test", 1 + 2, false)`,
		`@dump("test",(1 + 2),false)`,
	)
}

func TestIfDirString(t *testing.T) {
	t.Run("regular @if", func(t *testing.T) {
		inp := "@if(true)1@end"
		testDirectiveString[*ast.IfDir](t, inp, inp)
	})

	t.Run("with @else", func(t *testing.T) {
		inp := "@if(true)1@else2@end"
		testDirectiveString[*ast.IfDir](t, inp, inp)
	})

	t.Run("with @elseif", func(t *testing.T) {
		inp := "@if(true)first@elseif(false)second@end"
		testDirectiveString[*ast.IfDir](t, inp, inp)
	})

	t.Run("with @elseif and @else", func(t *testing.T) {
		inp := "@if(true)1@elseif(false)2@else3@end"
		testDirectiveString[*ast.IfDir](t, inp, inp)
	})
}

func TestForDirString(t *testing.T) {
	t.Run("regular @for", func(t *testing.T) {
		testDirectiveString[*ast.ForDir](
			t,
			"@for(i = 0; i < 10; i++){{ i }}@end",
			"@for((i = 0); (i < 10); (i++)){{ i }}@end",
		)
	})

	t.Run("with @else", func(t *testing.T) {
		testDirectiveString[*ast.ForDir](
			t,
			"@for(i = 0; i < 0; i++){{ i }}@elseEmpty@end",
			"@for((i = 0); (i < 0); (i++)){{ i }}@elseEmpty@end",
		)
	})

	t.Run("infinite @for", func(t *testing.T) {
		testDirectiveString[*ast.ForDir](t, "@for(;;)1@end", "@for(; ; )1@end")
	})
}

func TestEachDirString(t *testing.T) {
	t.Run("regular @each", func(t *testing.T) {
		inp := `@each(name in ["anna", "serhii"]){{ name }}@end`
		testDirectiveString[*ast.EachDir](t, inp, inp)
	})

	t.Run("with @else", func(t *testing.T) {
		inp := "@each(v in []){{ v }}@elseTest@end"
		testDirectiveString[*ast.EachDir](t, inp, inp)
	})
}

func TestProgramStringWithEmptyBlocks(t *testing.T) {
	cases := []struct {
		name   string
		inp    string
		expect string
	}{
		{
			name:   "empty if block",
			inp:    "@if(true)@end",
			expect: "@if(true)@end",
		},
		{
			name:   "empty if else block",
			inp:    "@if(true)a@else@end",
			expect: "@if(true)a@else@end",
		},
		{
			name:   "empty elseif block",
			inp:    "@if(true)a@elseif(false)@end",
			expect: "@if(true)a@elseif(false)@end",
		},
		{
			name:   "empty elseif else bodies",
			inp:    "@if(true)a@elseif(false)@else@end",
			expect: "@if(true)a@elseif(false)@else@end",
		},
		{
			name:   "empty if elseif else bodies",
			inp:    "@if(true)@elseif(false)@else@end",
			expect: "@if(true)@elseif(false)@else@end",
		},
		{
			name:   "empty each block",
			inp:    "@each(x in [1])@end",
			expect: "@each(x in [1])@end",
		},
		{
			name:   "empty each else",
			inp:    "@each(x in [1])a@else@end",
			expect: "@each(x in [1])a@else@end",
		},
		{
			name:   "empty for block",
			inp:    "@for(i=0;i<1;i++)@end",
			expect: "@for((i = 0); (i < 1); (i++))@end",
		},
		{
			name:   "empty for else block",
			inp:    "@for(i=0;i<1;i++)a@else@end",
			expect: "@for((i = 0); (i < 1); (i++))a@else@end",
		},
		{
			name:   "empty insert block",
			inp:    `@insert("a")@end`,
			expect: `@insert("a")@end`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := New(lexer.New(tc.inp), nil)
			prog := p.ParseProgram()

			if p.HasErrors() {
				t.Fatalf("parse errors for %q: %v", tc.inp, p.Errors()[0])
			}

			if got := prog.String(); got != tc.expect {
				t.Fatalf("Program.String() for %q = %q, expect %q", tc.inp, got, tc.expect)
			}
		})
	}
}

func TestBreakDirString(t *testing.T) {
	testChunkInLoopString[*ast.BreakDir](t, "@break", "@break")
}

func TestContinueDirString(t *testing.T) {
	testChunkInLoopString[*ast.ContinueDir](t, "@continue", "@continue")
}

func TestBreakifDirString(t *testing.T) {
	testChunkInLoopString[*ast.BreakifDir](t, "@breakif(true)", "@breakif(true)")
}

func TestContinueifDirString(t *testing.T) {
	testChunkInLoopString[*ast.ContinueifDir](t, "@continueif(false)", "@continueif(false)")
}
