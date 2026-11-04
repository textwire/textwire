package parser

import (
	"testing"

	"codeberg.org/textwire/textwire/v5/pkg/ast"
	"codeberg.org/textwire/textwire/v5/pkg/lexer"
)

func TestIdentExprString(t *testing.T) {
	inp := "{{ myName }}"
	expect := "myName"

	identExpr, err := parseEmbedded[*ast.IdentExpr](inp)
	if err != nil {
		t.Fatal(err)
	}

	if err := testString(identExpr, expect); err != nil {
		t.Fatal(err)
	}
}

func TestIntExprString(t *testing.T) {
	inp := "{{ 234 }}"
	expect := "234"

	intExpr, err := parseEmbedded[*ast.IntExpr](inp)
	if err != nil {
		t.Fatal(err)
	}

	if err := testString(intExpr, expect); err != nil {
		t.Fatal(err)
	}
}

func TestFloatExprString(t *testing.T) {
	inp := "{{ 2.34149 }}"
	expect := "2.34149"

	floatExpr, err := parseEmbedded[*ast.FloatExpr](inp)
	if err != nil {
		t.Fatal(err)
	}

	if err := testString(floatExpr, expect); err != nil {
		t.Fatal(err)
	}
}

func TestNilExprString(t *testing.T) {
	inp := "{{ nil }}"
	expect := "nil"

	nilExpr, err := parseEmbedded[*ast.NilExpr](inp)
	if err != nil {
		t.Fatal(err)
	}

	if err := testString(nilExpr, expect); err != nil {
		t.Fatal(err)
	}
}

func TestStrExprString(t *testing.T) {
	inp := "{{ 'Serhii' }}"
	expect := `"Serhii"`

	strExpr, err := parseEmbedded[*ast.StrExpr](inp)
	if err != nil {
		t.Fatal(err)
	}

	if err := testString(strExpr, expect); err != nil {
		t.Fatal(err)
	}
}

func TestBoolExprString(t *testing.T) {
	inp := "{{ true }}"
	expect := "true"

	boolExpr, err := parseEmbedded[*ast.BoolExpr](inp)
	if err != nil {
		t.Fatal(err)
	}

	if err := testString(boolExpr, expect); err != nil {
		t.Fatal(err)
	}
}

func TestInfixExprString(t *testing.T) {
	inp := "{{ 3 / 2 }}"
	expect := "(3 / 2)"

	infixExpr, err := parseEmbedded[*ast.InfixExpr](inp)
	if err != nil {
		t.Fatal(err)
	}

	if err := testString(infixExpr, expect); err != nil {
		t.Fatal(err)
	}
}

func TestPrefixExprString(t *testing.T) {
	inp := "{{ -obj.test }}"
	expect := "(-(obj.test))"

	prefixExpr, err := parseEmbedded[*ast.PrefixExpr](inp)
	if err != nil {
		t.Fatal(err)
	}

	if err := testString(prefixExpr, expect); err != nil {
		t.Fatal(err)
	}
}

func TestTernaryExprString(t *testing.T) {
	inp := "{{ true ? 1 : 0 }}"
	expect := "(true ? 1 : 0)"

	ternaryExpr, err := parseEmbedded[*ast.TernaryExpr](inp)
	if err != nil {
		t.Fatal(err)
	}

	if err := testString(ternaryExpr, expect); err != nil {
		t.Fatal(err)
	}
}

func TestArrExprString(t *testing.T) {
	inp := "{{ [11, 234,] }}"
	expect := "[11, 234]"

	arrExpr, err := parseEmbedded[*ast.ArrExpr](inp)
	if err != nil {
		t.Fatal(err)
	}

	if err := testString(arrExpr, expect); err != nil {
		t.Fatal(err)
	}
}

func TestObjExprString(t *testing.T) {
	inp := `{{ {"father": {"name": "John"}} }}`
	expect := `{"father": {"name": "John"}}`

	objExpr, err := parseEmbedded[*ast.ObjExpr](inp)
	if err != nil {
		t.Fatal(err)
	}

	if err := testString(objExpr, expect); err != nil {
		t.Fatal(err)
	}
}

func TestDotExprString(t *testing.T) {
	inp := "{{ person.father.name }}"
	expect := "((person.father).name)"

	dotExpr, err := parseEmbedded[*ast.DotExpr](inp)
	if err != nil {
		t.Fatal(err)
	}

	if err := testString(dotExpr, expect); err != nil {
		t.Fatal(err)
	}
}

func TestIndexExprString(t *testing.T) {
	inp := "{{ arr[1 + 2][2] }}"
	expect := "((arr[(1 + 2)])[2])"

	indexExpr, err := parseEmbedded[*ast.IndexExpr](inp)
	if err != nil {
		t.Fatal(err)
	}

	if err := testString(indexExpr, expect); err != nil {
		t.Fatal(err)
	}
}

func TestCallExprString(t *testing.T) {
	inp := "{{ name.len() }}"
	expect := "(name.len())"

	callExpr, err := parseEmbedded[*ast.CallExpr](inp)
	if err != nil {
		t.Fatal(err)
	}

	if err := testString(callExpr, expect); err != nil {
		t.Fatal(err)
	}
}

func TestGlobalCallExprString(t *testing.T) {
	inp := "{{ defined(name) }}"
	expect := "(defined(name))"

	globalCallExpr, err := parseEmbedded[*ast.GlobalCallExpr](inp)
	if err != nil {
		t.Fatal(err)
	}

	if err := testString(globalCallExpr, expect); err != nil {
		t.Fatal(err)
	}
}

func TestIncStmtString(t *testing.T) {
	inp := "{{ i++ }}"
	expect := "(i++)"

	incStmt, err := parseEmbedded[*ast.IncStmt](inp)
	if err != nil {
		t.Fatal(err)
	}

	if err := testString(incStmt, expect); err != nil {
		t.Fatal(err)
	}
}

func TestDecStmtString(t *testing.T) {
	inp := "{{ i-- }}"
	expect := "(i--)"

	decStmt, err := parseEmbedded[*ast.DecStmt](inp)
	if err != nil {
		t.Fatal(err)
	}

	if err := testString(decStmt, expect); err != nil {
		t.Fatal(err)
	}
}

func TestAssignStmtString(t *testing.T) {
	cases := []struct {
		inp    string
		expect string
	}{
		{inp: `{{ name = "Anna" }}`, expect: `(name = "Anna")`},
		{inp: `{{ myAge = 34 }}`, expect: `(myAge = 34)`},
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
			assignStmt, err := parseEmbedded[*ast.AssignStmt](tc.inp)
			if err != nil {
				t.Fatal(err)
			}

			if err := testString(assignStmt, tc.expect); err != nil {
				t.Fatal(err)
			}
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
	inp := "{!! 'x' !!}"
	expect := `{!! "x" !!}`

	embedded, err := parseDirective[*ast.Embedded](inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	if err := testString(embedded, expect); err != nil {
		t.Fatal(err)
	}
}

func TestTextString(t *testing.T) {
	inp := "<div><span>Hello</span></div>"

	text, err := parseDirective[*ast.Text](inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	if err := testString(text, inp); err != nil {
		t.Fatal(err)
	}
}

func TestUseDirString(t *testing.T) {
	inp := `@use("main")`

	useDir, err := parseDirective[*ast.UseDir](inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	if err := testString(useDir, inp); err != nil {
		t.Fatal(err)
	}
}

func TestReserveDirString(t *testing.T) {
	t.Run("without fallback", func(t *testing.T) {
		inp := `@reserve("content")`

		reserveDir, err := parseDirective[*ast.ReserveDir](inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}

		if err := testString(reserveDir, inp); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("with fallback", func(t *testing.T) {
		inp := `@reserve("content", 'fallback')`
		expect := `@reserve("content", "fallback")`

		reserveDir, err := parseDirective[*ast.ReserveDir](inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}

		if err := testString(reserveDir, expect); err != nil {
			t.Fatal(err)
		}
	})
}

func TestInsertDirString(t *testing.T) {
	t.Run("with block", func(t *testing.T) {
		inp := `@insert("content")<h1>Some content</h1>@end`

		insertDir, err := parseDirective[*ast.InsertDir](inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}

		if err := testString(insertDir, inp); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("with argument", func(t *testing.T) {
		inp := "@insert('content', 'Some content')"
		expect := `@insert("content", "Some content")`

		insertDir, err := parseDirective[*ast.InsertDir](inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}

		if err := testString(insertDir, expect); err != nil {
			t.Fatal(err)
		}
	})
}

func TestCompDirString(t *testing.T) {
	t.Run("without @pass", func(t *testing.T) {
		inp := "@component('components/book-card', { c: card })@end"
		expect := `@component("components/book-card", {"c": card})@end`

		compDir, err := parseDirective[*ast.CompDir](inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}

		if err := testString(compDir, expect); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("with default @pass", func(t *testing.T) {
		inp := `@component("components/book-card")<h1>Header</h1>@end`

		compDir, err := parseDirective[*ast.CompDir](inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}

		if err := testString(compDir, inp); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("with 2 @pass directives", func(t *testing.T) {
		inp := `@component("x")@pass("header")<h1>Header</h1>@end@pass("footer")<footer>Footer</footer>@end@end`

		compDir, err := parseDirective[*ast.CompDir](inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}

		if err := testString(compDir, inp); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("with @passif", func(t *testing.T) {
		inp := `@component("user")@passif(false, "name")Test2@end@end`

		compDir, err := parseDirective[*ast.CompDir](inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}

		if err := testString(compDir, inp); err != nil {
			t.Fatal(err)
		}
	})
}

func TestSlotDirString(t *testing.T) {
	t.Run("named slot", func(t *testing.T) {
		inp := "@slot('header')"
		expect := `@slot("header")`

		slotDir, err := parseDirective[*ast.SlotDir](inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}

		if err := testString(slotDir, expect); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("default slot", func(t *testing.T) {
		inp := "@slot"

		slotDir, err := parseDirective[*ast.SlotDir](inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}

		if err := testString(slotDir, inp); err != nil {
			t.Fatal(err)
		}
	})
}

func TestDumpDirString(t *testing.T) {
	inp := `@dump("test", 1 + 2, false)`
	expect := `@dump("test",(1 + 2),false)`

	dumpDir, err := parseDirective[*ast.DumpDir](inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	if err := testString(dumpDir, expect); err != nil {
		t.Fatal(err)
	}
}

func TestIfDirString(t *testing.T) {
	t.Run("regular @if", func(t *testing.T) {
		inp := "@if(true)1@end"

		ifDir, err := parseDirective[*ast.IfDir](inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}

		if err := testString(ifDir, inp); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("with @else", func(t *testing.T) {
		inp := "@if(true)1@else2@end"

		ifDir, err := parseDirective[*ast.IfDir](inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}

		if err := testString(ifDir, inp); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("with @elseif", func(t *testing.T) {
		inp := "@if(true)first@elseif(false)second@end"

		ifDir, err := parseDirective[*ast.IfDir](inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}

		if err := testString(ifDir, inp); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("with @elseif and @else", func(t *testing.T) {
		inp := "@if(true)1@elseif(false)2@else3@end"

		ifDir, err := parseDirective[*ast.IfDir](inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}

		if err := testString(ifDir, inp); err != nil {
			t.Fatal(err)
		}
	})
}

func TestForDirString(t *testing.T) {
	t.Run("regular @for", func(t *testing.T) {
		inp := "@for(i = 0; i < 10; i++){{ i }}@end"
		expect := "@for((i = 0); (i < 10); (i++)){{ i }}@end"

		forDir, err := parseDirective[*ast.ForDir](inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}

		if err := testString(forDir, expect); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("with @else", func(t *testing.T) {
		inp := "@for(i = 0; i < 0; i++){{ i }}@elseEmpty@end"
		expect := "@for((i = 0); (i < 0); (i++)){{ i }}@elseEmpty@end"

		forDir, err := parseDirective[*ast.ForDir](inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}

		if err := testString(forDir, expect); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("infinite @for", func(t *testing.T) {
		inp := "@for(;;)1@end"
		expect := "@for(; ; )1@end"

		forDir, err := parseDirective[*ast.ForDir](inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}

		if err := testString(forDir, expect); err != nil {
			t.Fatal(err)
		}
	})
}

func TestEachDirString(t *testing.T) {
	t.Run("regular @each", func(t *testing.T) {
		inp := `@each(name in ["anna", "serhii"]){{ name }}@end`

		eachDir, err := parseDirective[*ast.EachDir](inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}

		if err := testString(eachDir, inp); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("with @else", func(t *testing.T) {
		inp := "@each(v in []){{ v }}@elseTest@end"

		eachDir, err := parseDirective[*ast.EachDir](inp, defaultParseOpts)
		if err != nil {
			t.Fatal(err)
		}

		if err := testString(eachDir, inp); err != nil {
			t.Fatal(err)
		}
	})
}

func TestBreakDirString(t *testing.T) {
	inp := "@break"

	breakDir, err := parseChunkInLoop[*ast.BreakDir](inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	if err := testString(breakDir, inp); err != nil {
		t.Fatal(err)
	}
}

func TestContinueDirString(t *testing.T) {
	inp := "@continue"

	continueDir, err := parseChunkInLoop[*ast.ContinueDir](inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	if err := testString(continueDir, inp); err != nil {
		t.Fatal(err)
	}
}

func TestBreakifDirString(t *testing.T) {
	inp := "@breakif(true)"

	breakifDir, err := parseChunkInLoop[*ast.BreakifDir](inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	if err := testString(breakifDir, inp); err != nil {
		t.Fatal(err)
	}
}

func TestContinueifDirString(t *testing.T) {
	inp := "@continueif(false)"

	continueifDir, err := parseChunkInLoop[*ast.ContinueifDir](inp, defaultParseOpts)
	if err != nil {
		t.Fatal(err)
	}

	if err := testString(continueifDir, inp); err != nil {
		t.Fatal(err)
	}
}
