package evaluator

import (
	"testing"

	"codeberg.org/textwire/textwire/v6/pkg/fail"
	"codeberg.org/textwire/textwire/v6/pkg/value"
)

func TestEvalCircularReference(t *testing.T) {
	cases := []struct {
		name  string
		inp   string
		right string
		left  string
	}{
		{"array element assigned to its own array", `{{ a = [1]; a[0] = a; a }}`, "a", "a[0]"},
		{
			"object property assigned to its own object",
			`{{ o = {}; o.self = o; o }}`,
			"o",
			"o.self",
		},
		{
			"array element assigned to an array containing it",
			`{{ a = [1]; b = [a]; a[0] = b; a }}`,
			"b",
			"a[0]",
		},
		{
			"array element assigned to a nested reference of itself",
			`{{ a = [1]; b = [a]; a[0] = b[0]; a }}`,
			"b[0]",
			"a[0]",
		},
		{
			"object property assigned to an object containing it",
			`{{ o = {}; p = {o: o}; o.p = p; o }}`,
			"p",
			"o.p",
		},
		{
			"mutually referencing objects",
			`{{ o = {}; p = {}; o.p = p; p.o = o; o }}`,
			"o",
			"p.o",
		},
		{
			"nested array element assigned to the root array",
			`{{ a = [[1]]; a[0][0] = a; a }}`,
			"a",
			"a[0][0]",
		},
		{
			"object bracket property assigned to its own object",
			`{{ o = {}; o["self"] = o; o }}`,
			"o",
			`o["self"]`,
		},
		{
			"object bracket property assigned to a nested object chain",
			`{{ o = {}; p = {o: o}; o["p"] = p; o }}`,
			"p",
			`o["p"]`,
		},
	}

	for _, tc := range cases {
		evaluated, failure := testEval(tc.inp, nil)
		if failure != nil {
			t.Fatalf("case %s: evaluation failed: %s", tc.name, failure)
		}

		err, ok := evaluated.(*value.Error)
		if !ok {
			t.Fatalf("case %s: expected circular reference error, got %q", tc.name, evaluated)
		}

		expect := fail.New(
			nil,
			"/path/to/file",
			fail.OriginEval,
			fail.ErrCircularReference,
			tc.right,
			tc.left,
		)

		if err.String() != expect.String() {
			t.Fatalf(
				"case %s: error message must be:\n%q\ngot:\n%q",
				tc.name,
				expect,
				err,
			)
		}
	}
}

func TestEvalAssignAllowsSharedValues(t *testing.T) {
	cases := []struct {
		name   string
		inp    string
		expect string
	}{
		{"self assignment of an array element", `{{ a = [1]; a[0] = a[0]; a }}`, "1"},
		{"array element assigned a different array", `{{ a = [1]; b = [2]; a[0] = b; a }}`, "2"},
		{
			"array element assigned a shared nested array",
			`{{ a = [1]; b = [a]; b[0] = a; b }}`,
			"1",
		},
		{
			"object property assigned a separate object",
			`{{ o = {}; p = {}; o.p = p; o }}`,
			"{p: {}}",
		},
		{"array wrapped in a new array", `{{ a = [1]; a = [a]; a[0][0] }}`, "1"},
		{"array assigned to itself", `{{ a = [1]; a = a; a }}`, "1"},
		{
			"shared array mutated through a nested reference",
			`{{ a = [1]; b = [a]; b[0][0] = 2; b }}`,
			"2",
		},
		{
			"slice result assigned back into the original array",
			`{{ a = [1]; b = a.slice(0); b[0] = a; b }}`,
			"1",
		},
		{
			"same array shared twice in one array",
			`{{ a = [1]; b = [a, a]; b[0][0] = 2; b[1][0] }}`,
			"2",
		},
		{
			"same array shared across object properties",
			`{{ a = [1]; o = {}; o.p = a; o.q = a; o.p[0] = 2; o.q[0] }}`,
			"2",
		},
		{
			"same array shared across two arrays",
			`{{ a = [1]; b = [a]; c = [a]; b[0][0] = 2; c[0][0] }}`,
			"2",
		},
		{
			"replacing one shared element does not affect the other",
			`{{ a = [1]; b = [a, a]; b[0] = [2]; b[1][0] }}`,
			"1",
		},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, nil, tc.name)
	}
}

func TestObjCamelFuncRejectsCycle(t *testing.T) {
	obj := value.NewObj(nil)
	obj.Pairs["self"] = obj

	if _, err := objCamelFunc(obj); err == nil {
		t.Fatal("expected circular reference error")
	}
}
