package evaluator

import (
	"testing"

	"codeberg.org/textwire/textwire/v5/pkg/fail"
	"codeberg.org/textwire/textwire/v5/pkg/value"
)

func TestEvalCircularReference(t *testing.T) {
	cases := []struct {
		id    uint
		inp   string
		right string
		left  string
	}{
		{10, `{{ a = [1]; a[0] = a; a }}`, "a", "a[0]"},
		{20, `{{ o = {}; o.self = o; o }}`, "o", "o.self"},
		{30, `{{ a = [1]; b = [a]; a[0] = b; a }}`, "b", "a[0]"},
		{40, `{{ a = [1]; b = [a]; a[0] = b[0]; a }}`, "b[0]", "a[0]"},
		{50, `{{ o = {}; p = {o: o}; o.p = p; o }}`, "p", "o.p"},
		{60, `{{ o = {}; p = {}; o.p = p; p.o = o; o }}`, "o", "p.o"},
		{70, `{{ a = [[1]]; a[0][0] = a; a }}`, "a", "a[0][0]"},
		{80, `{{ o = {}; o["self"] = o; o }}`, "o", `o["self"]`},
		{90, `{{ o = {}; p = {o: o}; o["p"] = p; o }}`, "p", `o["p"]`},
	}

	for _, tc := range cases {
		evaluated, failure := testEval(tc.inp, nil)
		if failure != nil {
			t.Fatalf("Case: %d. evaluation failed: %s", tc.id, failure)
		}

		err, ok := evaluated.(*value.Error)
		if !ok {
			t.Fatalf("Case: %d. expected circular reference error, got %q", tc.id, evaluated)
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
				"Case: %d. Error message must be:\n%q\ngot:\n%q",
				tc.id,
				expect,
				err,
			)
		}
	}
}

func TestEvalAssignAllowsSharedValues(t *testing.T) {
	cases := []struct {
		id     uint
		inp    string
		expect string
	}{
		{10, `{{ a = [1]; a[0] = a[0]; a }}`, "1"},
		{20, `{{ a = [1]; b = [2]; a[0] = b; a }}`, "2"},
		{30, `{{ a = [1]; b = [a]; b[0] = a; b }}`, "1"},
		{40, `{{ o = {}; p = {}; o.p = p; o }}`, "{p: {}}"},
		{50, `{{ a = [1]; a = [a]; a[0][0] }}`, "1"},
		{60, `{{ a = [1]; a = a; a }}`, "1"},
		{70, `{{ a = [1]; b = [a]; b[0][0] = 2; b }}`, "2"},
		{80, `{{ a = [1]; b = a.slice(0); b[0] = a; b }}`, "1"},
		{90, `{{ a = [1]; b = [a, a]; b[0][0] = 2; b[1][0] }}`, "2"},
		{100, `{{ a = [1]; o = {}; o.p = a; o.q = a; o.p[0] = 2; o.q[0] }}`, "2"},
		{110, `{{ a = [1]; b = [a]; c = [a]; b[0][0] = 2; c[0][0] }}`, "2"},
		{120, `{{ a = [1]; b = [a, a]; b[0] = [2]; b[1][0] }}`, "1"},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, nil, tc.id)
	}
}

func TestObjCamelFuncRejectsCycle(t *testing.T) {
	obj := value.NewObj(nil)
	obj.Pairs["self"] = obj

	if _, err := objCamelFunc(obj); err == nil {
		t.Fatal("expected circular reference error")
	}
}
