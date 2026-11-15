package value

import (
	"fmt"
	"math"
	"reflect"
	"strings"
	"testing"
	"time"

	"codeberg.org/textwire/textwire/v6/pkg/fail"
)

// Named types used to verify that NativeToValue looks at the underlying kind.
type (
	testRole  string
	testCents int
	testLevel uint8
	testRatio float64
	testFlag  bool
)

type testUser struct {
	Name string
	Role testRole
}

type testModel struct {
	Name     string
	Callback func()
	Ch       chan int
	Data     [3]int
}

type testUintModel struct {
	Name string
	Big  uint64
}

func intsArray(vals ...int64) *Arr {
	elems := make([]Literal, len(vals))
	for i, v := range vals {
		elems[i] = &Int{Val: v}
	}
	return &Arr{Elements: elems}
}

func TestNativeToValue_SupportedValues(t *testing.T) {
	role := testRole("admin")
	createdAt := time.Date(2026, 10, 26, 12, 30, 0, 0, time.UTC)

	tests := []struct {
		name   string
		inp    any
		expect Literal
	}{
		// named primitives
		{"named string", testRole("editor"), &Str{Val: "editor"}},
		{"named int", testCents(42), &Int{Val: 42}},
		{"named uint", testLevel(7), &Int{Val: 7}},
		{"named float", testRatio(1.5), &Float{Val: 1.5}},
		{"named bool true", testFlag(true), &Bool{Val: true}},
		{"named bool false", testFlag(false), &Bool{Val: false}},

		// stdlib types
		{"duration", time.Duration(5), &Int{Val: 5}},
		{"time", createdAt, &Str{Val: "2026-10-26 12:30:00"}},

		// pointers and composites
		{"pointer to named string", &role, &Str{Val: "admin"}},
		{"fixed-size array", [3]int{1, 2, 3}, intsArray(1, 2, 3)},
		{
			"map with integer keys",
			map[int]string{1: "one", 2: "two"},
			NewObj(map[string]Literal{
				"1": &Str{Val: "one"},
				"2": &Str{Val: "two"},
			}),
		},
		{
			name: "struct with named field",
			inp:  testUser{Name: "Bob", Role: "editor"},
			expect: NewObj(map[string]Literal{
				"Name": &Str{Val: "Bob"},
				"Role": &Str{Val: "editor"},
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertValue(t, tt.inp, tt.expect)
		})
	}
}

func TestNativeToValue_UnsupportedValues(t *testing.T) {
	tests := []struct {
		name      string
		inp       any
		expectErr string
	}{
		{"func", func() {}, "unsupported value type 'func()'"},
		{"chan in", make(chan int), "unsupported value type 'chan int'"},
		{"chan string", make(chan string), "unsupported value type 'chan string'"},
		{"complex", complex(1, 2), "unsupported value type 'complex128'"},
		{"complex 128", complex128(2), "unsupported value type 'complex128'"},
		{"complex 64", complex64(2), "unsupported value type 'complex64'"},
		{"slice containing func", []any{func() {}}, "unsupported value type 'func()'"},
		{
			"map containing func",
			map[string]any{"drop": func() {}},
			"unsupported value type 'func()'",
		},
		{
			name: "struct with func field",
			inp: testModel{
				Name:     "Bob",
				Callback: func() {},
				Ch:       make(chan int),
				Data:     [3]int{1, 2, 3},
			},
			expectErr: "unsupported value type 'func()' for field 'Callback'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertErrorMessage(t, tt.inp, tt.expectErr)
		})
	}
}

func TestNativeToValue_UintOverflow(t *testing.T) {
	t.Run("max int64 fits", func(t *testing.T) {
		assertValue(t, uint64(math.MaxInt64), &Int{Val: math.MaxInt64})
	})

	t.Run("above max int64 errors", func(t *testing.T) {
		val := uint64(math.MaxInt64) + 1
		assertErrorMessage(t, val, fmt.Sprintf(fail.ErrUintOverflow, val))
	})

	t.Run("above max int64 in a field names the field", func(t *testing.T) {
		inp := testUintModel{Name: "Bob", Big: uint64(math.MaxUint64)}
		assertErrorMessage(t, inp, fmt.Sprintf(fail.ErrUintOverflowField, inp.Big, "Big"))
	})
}

func TestNativeToValue_CircularValues(t *testing.T) {
	t.Run("map cycle", func(t *testing.T) {
		inp := map[string]any{}
		inp["self"] = inp

		assertErrorMessage(t, inp, fmt.Sprintf(fail.ErrCircularValue, "map[string]interface {}"))
	})

	t.Run("slice cycle", func(t *testing.T) {
		inp := make([]any, 1)
		inp[0] = inp

		assertErrorMessage(t, inp, fmt.Sprintf(fail.ErrCircularValue, "[]interface {}"))
	})

	t.Run("struct pointer cycle", func(t *testing.T) {
		type node struct {
			Next *node
		}

		inp := &node{}
		inp.Next = inp

		got := NativeToValue(inp)
		errObj, ok := got.(*Error)
		if !ok {
			t.Fatalf("NativeToValue(%T) returned %T, expect *Error", inp, got)
		}

		if !strings.Contains(errObj.Err.Message(), "circular reference") {
			t.Fatalf("unexpected error: %q", errObj.Err.Message())
		}
	})

	t.Run("shared value without cycle", func(t *testing.T) {
		shared := map[string]any{"x": 1}
		inp := map[string]any{"a": shared, "b": shared}

		if got := NativeToValue(inp); got.Is(ERR_VAL) {
			t.Fatalf("unexpected error: %s", got)
		}
	})

	t.Run("parent pointer without cycle", func(t *testing.T) {
		type Comment struct {
			Parent *Comment
		}

		root := &Comment{}
		child := &Comment{Parent: root}

		if got := NativeToValue(child); got.Is(ERR_VAL) {
			t.Fatalf("unexpected error: %s", got)
		}
	})

	t.Run("pointer to first field without cycle", func(t *testing.T) {
		type Inner struct{ X int }
		type Outer struct {
			In  Inner
			Ref *Inner
		}

		outer := &Outer{}
		outer.Ref = &outer.In

		if got := NativeToValue(outer); got.Is(ERR_VAL) {
			t.Fatalf("unexpected error: %s", got)
		}
	})

	t.Run("sub-slice cycle", func(t *testing.T) {
		inp := make([]any, 2)
		inp[0] = inp[:1]

		if got := NativeToValue(inp); !got.Is(ERR_VAL) {
			t.Fatalf("expected circular reference error, got %s", got)
		}
	})

	t.Run("sub-slice of ancestor without cycle", func(t *testing.T) {
		inp := make([]any, 2)
		inp[0] = 1
		inp[1] = inp[:1]

		if got := NativeToValue(inp); got.Is(ERR_VAL) {
			t.Fatalf("unexpected error: %s", got)
		}
	})
}

func TestRefKeyFor(t *testing.T) {
	t.Run("slice key holds type, pointer and length", func(t *testing.T) {
		inp := make([]any, 2)

		key, ok := refKeyFor(reflect.ValueOf(inp))
		if !ok {
			t.Fatal("expected slice to be trackable")
		}

		expect := refKey{
			typ: reflect.TypeFor[[]any](),
			ptr: reflect.ValueOf(inp).Pointer(),
			len: 2,
		}

		if key != expect {
			t.Fatalf("expected %+v, got %+v", expect, key)
		}
	})

	t.Run("sub-slice key differs from its parent slice", func(t *testing.T) {
		inp := make([]any, 2)

		whole, _ := refKeyFor(reflect.ValueOf(inp))
		part, _ := refKeyFor(reflect.ValueOf(inp[:1]))
		if whole == part {
			t.Fatalf("expected different keys, got %+v", whole)
		}
	})

	t.Run("map key holds type and pointer", func(t *testing.T) {
		inp := map[string]any{}

		key, ok := refKeyFor(reflect.ValueOf(inp))
		if !ok {
			t.Fatal("expected map to be trackable")
		}

		expect := refKey{
			typ: reflect.TypeFor[map[string]any](),
			ptr: reflect.ValueOf(inp).Pointer(),
		}

		if key != expect {
			t.Fatalf("expected %+v, got %+v", expect, key)
		}
	})

	t.Run("non-nil pointer is trackable", func(t *testing.T) {
		inp := &struct{}{}

		if _, ok := refKeyFor(reflect.ValueOf(inp)); !ok {
			t.Fatal("expected pointer to be trackable")
		}
	})

	t.Run("nil pointer is not trackable", func(t *testing.T) {
		var inp *struct{}

		if _, ok := refKeyFor(reflect.ValueOf(inp)); ok {
			t.Fatal("expected nil pointer not to be trackable")
		}
	})

	t.Run("int is not trackable", func(t *testing.T) {
		if _, ok := refKeyFor(reflect.ValueOf(1)); ok {
			t.Fatal("expected int not to be trackable")
		}
	})

	t.Run("pointer key holds type and address", func(t *testing.T) {
		inp := &struct{ X int }{}

		key, ok := refKeyFor(reflect.ValueOf(inp))
		if !ok {
			t.Fatal("expected pointer to be trackable")
		}

		expect := refKey{
			typ: reflect.TypeFor[*struct{ X int }](),
			ptr: reflect.ValueOf(inp).Pointer(),
		}

		if key != expect {
			t.Fatalf("expected %+v, got %+v", expect, key)
		}
	})

	t.Run("different types at same address have different keys", func(t *testing.T) {
		type Inner struct{ X int }
		type Outer struct{ In Inner }

		outer := &Outer{}
		inner := &outer.In

		outerValue := reflect.ValueOf(outer)
		innerValue := reflect.ValueOf(inner)
		if outerValue.Pointer() != innerValue.Pointer() {
			t.Fatal("expected Outer and its first field to share an address")
		}

		keyOuter, ok := refKeyFor(outerValue)
		if !ok {
			t.Fatal("expected Outer pointer to be trackable")
		}

		keyInner, ok := refKeyFor(innerValue)
		if !ok {
			t.Fatal("expected Inner pointer to be trackable")
		}

		if keyOuter == keyInner {
			t.Fatal("expected different keys for Outer and Inner pointers at same address")
		}
	})
}

func TestIsControlDir(t *testing.T) {
	t.Run("break and continue are control directives", func(t *testing.T) {
		if !IsControlDir(&Break{}) || !IsControlDir(&Continue{}) {
			t.Fatal("expected break and continue to be control directives")
		}
	})

	t.Run("text and blocks are not control directives", func(t *testing.T) {
		if IsControlDir(&Text{Val: "a"}) || IsControlDir(NewBlock(0)) {
			t.Fatal("expected text and block not to be control directives")
		}
	})
}

func TestHasBreak(t *testing.T) {
	t.Run("break value has break", func(t *testing.T) {
		if !HasBreak(&Break{}) {
			t.Fatal("expected break value to have break")
		}
	})

	t.Run("text value has no break", func(t *testing.T) {
		if HasBreak(&Text{Val: "a"}) {
			t.Fatal("expected text value to have no break")
		}
	})

	t.Run("block with nested break has break", func(t *testing.T) {
		inner := NewBlock(1)
		inner.Append(&Break{})

		block := NewBlock(1)
		block.Append(inner)

		if !HasBreak(block) {
			t.Fatal("expected nested break to be found")
		}
	})

	t.Run("block without break has no break", func(t *testing.T) {
		block := NewBlock(1)
		block.Append(&Text{Val: "a"})

		if HasBreak(block) {
			t.Fatal("expected block without break to have no break")
		}
	})
}

func TestHasContinue(t *testing.T) {
	t.Run("continue value has continue", func(t *testing.T) {
		if !HasContinue(&Continue{}) {
			t.Fatal("expected continue value to have continue")
		}
	})

	t.Run("text value has no continue", func(t *testing.T) {
		if HasContinue(&Text{Val: "a"}) {
			t.Fatal("expected text value to have no continue")
		}
	})

	t.Run("block with nested continue has continue", func(t *testing.T) {
		inner := NewBlock(1)
		inner.Append(&Continue{})

		block := NewBlock(1)
		block.Append(inner)

		if !HasContinue(block) {
			t.Fatal("expected nested continue to be found")
		}
	})

	t.Run("block without continue has no continue", func(t *testing.T) {
		block := NewBlock(1)
		block.Append(&Text{Val: "a"})

		if HasContinue(block) {
			t.Fatal("expected block without continue to have no continue")
		}
	})
}
