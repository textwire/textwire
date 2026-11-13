package value

import (
	"fmt"
	"testing"

	"codeberg.org/textwire/textwire/v6/pkg/fail"
)

func TestReferences(t *testing.T) {
	t.Run("direct array reference", func(t *testing.T) {
		target := &Arr{}
		root := &Arr{Elements: []Literal{target}}

		if !References(root, target) {
			t.Fatal("expected target to be referenced")
		}
	})

	t.Run("indirect object reference", func(t *testing.T) {
		target := &Obj{}
		mid := &Arr{Elements: []Literal{target}}
		root := &Obj{Pairs: map[string]Literal{"mid": mid}}

		if !References(root, target) {
			t.Fatal("expected target to be referenced")
		}
	})

	t.Run("unrelated values", func(t *testing.T) {
		target := &Arr{}
		root := &Arr{Elements: []Literal{&Arr{}}}

		if References(root, target) {
			t.Fatal("expected target not to be referenced")
		}
	})

	t.Run("existing cycle without target", func(t *testing.T) {
		cycle := &Arr{}
		cycle.Elements = []Literal{cycle}

		if References(cycle, &Arr{}) {
			t.Fatal("expected target not to be referenced")
		}
	})
}

func TestHasCycle(t *testing.T) {
	t.Run("no cycle", func(t *testing.T) {
		root := &Arr{Elements: []Literal{
			&Obj{Pairs: map[string]Literal{"x": &Int{Val: 1}}},
		}}

		if HasCycle(root) {
			t.Fatal("expected no cycle")
		}
	})

	t.Run("direct array cycle", func(t *testing.T) {
		root := &Arr{}
		root.Elements = []Literal{root}

		if !HasCycle(root) {
			t.Fatal("expected cycle")
		}
	})

	t.Run("indirect cycle", func(t *testing.T) {
		first := &Arr{}
		second := &Obj{Pairs: map[string]Literal{}}
		first.Elements = []Literal{second}
		second.Pairs["first"] = first

		if !HasCycle(first) {
			t.Fatal("expected cycle")
		}
	})

	t.Run("shared value without cycle", func(t *testing.T) {
		shared := &Arr{Elements: []Literal{&Int{Val: 1}}}
		root := &Arr{Elements: []Literal{shared, shared}}

		if HasCycle(root) {
			t.Fatal("expected no cycle")
		}
	})
}

func TestCycleSerialization(t *testing.T) {
	t.Run("array cycle", func(t *testing.T) {
		cycle := &Arr{}
		cycle.Elements = []Literal{cycle}

		if got := cycle.String(); got != circularReference {
			t.Fatalf("expected %q, got %q", circularReference, got)
		}

		if got := cycle.Dump(0); got != circularReference {
			t.Fatalf("expected %q, got %q", circularReference, got)
		}

		_, err := cycle.JSON()
		if err == nil {
			t.Fatal("expected JSON error")
		}

		expect := fmt.Sprintf(fail.ErrCircularValue, ARR_VAL)
		if err.Error() != expect {
			t.Fatalf("expected %q, got %q", expect, err.Error())
		}

		if got := cycle.Native(); got != nil {
			t.Fatalf("expected nil, got %#v", got)
		}
	})

	t.Run("object cycle", func(t *testing.T) {
		cycle := &Obj{Pairs: map[string]Literal{}}
		cycle.Pairs["self"] = cycle

		if got := cycle.String(); got != circularReference {
			t.Fatalf("expected %q, got %q", circularReference, got)
		}

		if got := cycle.Dump(0); got != circularReference {
			t.Fatalf("expected %q, got %q", circularReference, got)
		}

		_, err := cycle.JSON()
		if err == nil {
			t.Fatal("expected JSON error")
		}

		expect := fmt.Sprintf(fail.ErrCircularValue, OBJ_VAL)
		if err.Error() != expect {
			t.Fatalf("expected %q, got %q", expect, err.Error())
		}

		if got := cycle.Native(); got != nil {
			t.Fatalf("expected nil, got %#v", got)
		}
	})

	t.Run("indirect cycle", func(t *testing.T) {
		first := &Arr{}
		second := &Obj{Pairs: map[string]Literal{}}
		first.Elements = []Literal{second}
		second.Pairs["first"] = first

		if got := first.String(); got != circularReference {
			t.Fatalf("expected %q, got %q", circularReference, got)
		}
	})
}
