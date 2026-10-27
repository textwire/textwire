package value

import (
	"fmt"
	"math"
	"reflect"
	"testing"
	"time"

	"codeberg.org/textwire/textwire/v5/pkg/fail"
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

func assertValue(t *testing.T, input any, want Literal) {
	t.Helper()

	got := NativeToValue(input)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("NativeToValue(%T)\n got: %#v\nwant: %#v", input, got, want)
	}
}

func assertErrorMessage(t *testing.T, input any, want string) {
	t.Helper()

	got := NativeToValue(input)

	errObj, ok := got.(*Error)
	if !ok {
		t.Fatalf("NativeToValue(%T) returned %T, want *Error", input, got)
	}

	if msg := errObj.Err.Message(); msg != want {
		t.Fatalf("error message\n got: %q\nwant: %q", msg, want)
	}
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
		name  string
		input any
		want  Literal
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
			name:  "struct with named field",
			input: testUser{Name: "Bob", Role: "editor"},
			want: NewObj(map[string]Literal{
				"Name": &Str{Val: "Bob"},
				"Role": &Str{Val: "editor"},
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertValue(t, tt.input, tt.want)
		})
	}
}

func TestNativeToValue_UnsupportedValues(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		wantErr string
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
			input: testModel{
				Name:     "Bob",
				Callback: func() {},
				Ch:       make(chan int),
				Data:     [3]int{1, 2, 3},
			},
			wantErr: "unsupported value type 'func()' for field 'Callback'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertErrorMessage(t, tt.input, tt.wantErr)
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
		input := testUintModel{Name: "Bob", Big: uint64(math.MaxUint64)}
		assertErrorMessage(t, input, fmt.Sprintf(fail.ErrUintOverflowField, input.Big, "Big"))
	})
}
