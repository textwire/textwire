package value

import "testing"

func TestNewFrozenGlobal(t *testing.T) {
	cases := []struct {
		name      string
		key       string
		val       any
		expectVal Literal
	}{
		{
			name:      "preserves string values",
			key:       "name",
			val:       "Amy Adams 😀",
			expectVal: &Str{Val: "Amy Adams 😀"},
		},
		{
			name:      "preserves integer values",
			key:       "age",
			val:       25,
			expectVal: &Int{Val: 25},
		},
		{
			name:      "preserves negative integers",
			key:       "score",
			val:       -10,
			expectVal: &Int{Val: -10},
		},
		{
			name:      "preserves float values",
			key:       "height",
			val:       5.7,
			expectVal: &Float{Val: 5.7},
		},
		{
			name:      "preserves negative floats",
			key:       "temperature",
			val:       -2.5,
			expectVal: &Float{Val: -2.5},
		},
		{
			name:      "preserves true booleans",
			key:       "isActive",
			val:       true,
			expectVal: &Bool{Val: true},
		},
		{
			name:      "preserves false booleans",
			key:       "isComplete",
			val:       false,
			expectVal: &Bool{Val: false},
		},
		{
			name:      "preserves empty strings",
			key:       "empty",
			val:       "",
			expectVal: &Str{Val: ""},
		},
		{
			name:      "preserves zero integers",
			key:       "zero",
			val:       0,
			expectVal: &Int{Val: 0},
		},
		{
			name:      "preserves zero floats",
			key:       "zeroFloat",
			val:       0.0,
			expectVal: &Float{Val: 0.0},
		},
		{
			name:      "preserves strings with special characters",
			key:       "message",
			val:       "Hello, World! @#$%",
			expectVal: &Str{Val: "Hello, World! @#$%"},
		},
		{
			name:      "preserves large integers",
			key:       "bigNumber",
			val:       999999999,
			expectVal: &Int{Val: 999999999},
		},
		{
			name: "converts integer slices to arrays",
			key:  "numbers",
			val:  []int{1, 2, 3},
			expectVal: &Arr{Elements: []Literal{
				&Int{Val: 1},
				&Int{Val: 2},
				&Int{Val: 3},
			}},
		},
		{
			name: "converts string slices to arrays",
			key:  "names",
			val:  []string{"Alice", "Bob", "Charlie"},
			expectVal: &Arr{Elements: []Literal{
				&Str{Val: "Alice"},
				&Str{Val: "Bob"},
				&Str{Val: "Charlie"},
			}},
		},
		{
			name: "converts mixed type slices to arrays",
			key:  "mixed",
			val:  []any{"hello", 42, true},
			expectVal: &Arr{Elements: []Literal{
				&Str{Val: "hello"},
				&Int{Val: 42},
				&Bool{Val: true},
			}},
		},
		{
			name:      "converts empty slices to empty arrays",
			key:       "emptySlice",
			val:       []int{},
			expectVal: &Arr{Elements: []Literal{}},
		},
		{
			name: "converts maps to objects",
			key:  "user",
			val:  map[string]any{"name": "John", "age": 30},
			expectVal: NewObj(map[string]Literal{
				"name": &Str{Val: "John"},
				"age":  &Int{Val: 30},
			}),
		},
		{
			name: "converts nested maps to nested objects",
			key:  "config",
			val: map[string]any{
				"database": map[string]any{
					"host": "localhost",
					"port": 5432,
				},
			},
			expectVal: NewObj(map[string]Literal{
				"database": NewObj(map[string]Literal{
					"host": &Str{Val: "localhost"},
					"port": &Int{Val: 5432},
				}),
			}),
		},
		{
			name: "converts objects containing slices",
			key:  "data",
			val: map[string]any{
				"tags":  []string{"go", "test", "unit"},
				"count": 3,
			},
			expectVal: NewObj(map[string]Literal{
				"tags": &Arr{Elements: []Literal{
					&Str{Val: "go"},
					&Str{Val: "test"},
					&Str{Val: "unit"},
				}},
				"count": &Int{Val: 3},
			}),
		},
		{
			name: "converts slices of objects",
			key:  "users",
			val: []map[string]any{
				{"name": "Alice", "age": 25},
				{"name": "Bob", "age": 30},
			},
			expectVal: &Arr{Elements: []Literal{
				NewObj(map[string]Literal{
					"name": &Str{Val: "Alice"},
					"age":  &Int{Val: 25},
				}),
				NewObj(map[string]Literal{
					"name": &Str{Val: "Bob"},
					"age":  &Int{Val: 30},
				}),
			}},
		},
		{
			name:      "converts nil slices to empty arrays",
			key:       "nilSlice",
			val:       []any(nil),
			expectVal: &Arr{Elements: []Literal{}},
		},
		{
			name:      "converts empty maps to empty objects",
			key:       "emptyObj",
			val:       map[string]any{},
			expectVal: NewObj(nil),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			global, failure := NewFrozenGlobal(map[string]any{tc.key: tc.val})
			if failure != nil {
				t.Fatalf("unexpected error: %s", failure)
			}

			val, ok := global.Pairs[tc.key]
			if !ok {
				t.Fatalf("global does not have key %s", tc.key)
			}

			if val.Type() != tc.expectVal.Type() {
				t.Fatalf("expected type %q, got %q", tc.expectVal.Type(), val.Type())
			}

			if val.String() != tc.expectVal.String() {
				t.Fatalf("expected %q, got %q", tc.expectVal, val)
			}

			assertDeepFrozen(t, val)
		})
	}

	t.Run("nil input results in empty global", func(t *testing.T) {
		global, failure := NewFrozenGlobal(nil)
		if failure != nil {
			t.Fatalf("unexpected error: %s", failure)
		}

		if global != emptyGlobal {
			t.Fatal("expected nil data to return the shared emptyGlobal")
		}
	})

	t.Run("unsupported values results in failure", func(t *testing.T) {
		_, failure := NewFrozenGlobal(map[string]any{"ch": make(chan int)})
		if failure == nil {
			t.Fatal("expected error for unsupported global value, got none")
		}
	})

	t.Run("NewScope contains readonly global", func(t *testing.T) {
		global, ok := NewScope(nil).Get("global")
		if !ok {
			t.Fatal("expected global in a fresh scope")
		}

		obj, ok := global.(*Obj)
		if !ok || !obj.IsFrozen() {
			t.Fatal("expected global to be a frozen object")
		}
	})

	t.Run("cycling arrays and nested objects are frozen", func(t *testing.T) {
		arr := &Arr{}
		obj := NewObj(map[string]Literal{"self": arr})
		arr.Elements = []Literal{arr, obj}

		Freeze(arr)

		if !arr.IsFrozen() {
			t.Fatal("expected array to be frozen")
		}

		if !obj.IsFrozen() {
			t.Fatal("expected nested object to be frozen")
		}
	})
}

func assertDeepFrozen(t *testing.T, v Literal) {
	t.Helper()

	if !isFrozen(v) {
		t.Fatalf("expected %s to be frozen", v.Type())
	}

	kids, _ := childValues(v)
	for _, kid := range kids {
		assertDeepFrozen(t, kid)
	}
}

func isFrozen(v Literal) bool {
	switch val := v.(type) {
	case *Obj:
		return val.IsFrozen()
	case *Arr:
		return val.IsFrozen()
	}

	return true
}
