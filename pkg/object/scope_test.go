package object

import (
	"testing"
)

func TestNewScopeFromMap(t *testing.T) {
	var float32val float32 = 5.731

	data := map[string]any{
		"title":    "Hello, World!",
		"n":        -1,
		"num":      int8(-2),
		"num1":     int16(-3),
		"num2":     int32(-4),
		"num3":     int64(-5),
		"un":       uint(1),
		"unum":     uint8(2),
		"unum1":    uint16(3),
		"unum2":    uint32(4),
		"unum3":    uint64(5),
		"bool1":    true,
		"bool2":    false,
		"height":   5.7,
		"weight":   float32val,
		"nothing":  nil,
		"ages":     []int{1, 2, 3},
		"ages64":   []int64{1, 2, 3},
		"ages32":   []int32{1, 2, 3},
		"ages16":   []int16{1, 2, 3},
		"ages8":    []int8{1, 2, 3},
		"nums":     []uint{1, 2, 3},
		"nums64":   []uint64{1, 2, 3},
		"nums32":   []uint32{1, 2, 3},
		"nums16":   []uint16{1, 2, 3},
		"nums8":    []uint8{1, 2, 3},
		"names":    []string{"John", "Jane"},
		"statuses": []bool{false, true},
		"rates64":  []float64{23.4, 56.7, 89.0},
		"values":   []any{23.4, 56.7, 89.0},
		"rates32":  []float32{float32val, float32val, float32val},
	}

	expect := map[string]Object{
		"title":   &String{Val: "Hello, World!"},
		"n":       &Integer{Val: -1},
		"num":     &Integer{Val: -2},
		"num1":    &Integer{Val: -3},
		"num2":    &Integer{Val: -4},
		"num3":    &Integer{Val: -5},
		"un":      &Integer{Val: 1},
		"unum":    &Integer{Val: 2},
		"unum1":   &Integer{Val: 3},
		"unum2":   &Integer{Val: 4},
		"unum3":   &Integer{Val: 5},
		"bool1":   &Boolean{Val: true},
		"bool2":   &Boolean{Val: false},
		"height":  &Float{Val: 5.7},
		"weight":  &Float{Val: float64(float32val)},
		"nothing": new(Nil),
		"ages": &Array{
			Elements: []Object{&Integer{Val: 1}, &Integer{Val: 2}, &Integer{Val: 3}},
		},
		"ages64": &Array{
			Elements: []Object{&Integer{Val: 1}, &Integer{Val: 2}, &Integer{Val: 3}},
		},
		"ages32": &Array{
			Elements: []Object{&Integer{Val: 1}, &Integer{Val: 2}, &Integer{Val: 3}},
		},
		"ages16": &Array{
			Elements: []Object{&Integer{Val: 1}, &Integer{Val: 2}, &Integer{Val: 3}},
		},
		"ages8": &Array{
			Elements: []Object{&Integer{Val: 1}, &Integer{Val: 2}, &Integer{Val: 3}},
		},
		"nums": &Array{
			Elements: []Object{&Integer{Val: 1}, &Integer{Val: 2}, &Integer{Val: 3}},
		},
		"nums64": &Array{
			Elements: []Object{&Integer{Val: 1}, &Integer{Val: 2}, &Integer{Val: 3}},
		},
		"nums32": &Array{
			Elements: []Object{&Integer{Val: 1}, &Integer{Val: 2}, &Integer{Val: 3}},
		},
		"nums16": &Array{
			Elements: []Object{&Integer{Val: 1}, &Integer{Val: 2}, &Integer{Val: 3}},
		},
		"nums8": &Array{
			Elements: []Object{&Integer{Val: 1}, &Integer{Val: 2}, &Integer{Val: 3}},
		},
		"names":    &Array{Elements: []Object{&String{Val: "John"}, &String{Val: "Jane"}}},
		"statuses": &Array{Elements: []Object{&Boolean{Val: false}, &Boolean{Val: true}}},
		"rates64": &Array{
			Elements: []Object{&Float{Val: 23.4}, &Float{Val: 56.7}, &Float{Val: 89.0}},
		},
		"values": &Array{
			Elements: []Object{&Float{Val: 23.4}, &Float{Val: 56.7}, &Float{Val: 89.0}},
		},
		"rates32": &Array{
			Elements: []Object{
				&Float{Val: float64(float32val)},
				&Float{Val: float64(float32val)},
				&Float{Val: float64(float32val)},
			},
		},
	}

	scope, err := NewScopeFromMap(data)
	if err != nil {
		t.Fatalf("returned an error: %s", err)
	}

	for key, val := range expect {
		obj, ok := scope.Get(key)
		if !ok {
			t.Fatalf("scope.Get(%s) returned !ok", key)
		}

		if obj.Type() != val.Type() {
			t.Fatalf("scope.Get(%s) returned %q, expected %q", key, obj.Type(), val.Type())
		}

		if obj.String() != val.String() {
			t.Fatalf("scope.Get(%s) returned %q, expected %q", key, obj, val)
		}
	}
}

func TestAddGlobalData(t *testing.T) {
	cases := []struct {
		name      string
		key       string
		val       any
		expectKey string
		expectVal Object
	}{
		{
			name:      "add string global identifier",
			key:       "name",
			val:       "Amy Adams",
			expectKey: "name",
			expectVal: &String{Val: "Amy Adams"},
		},
		{
			name:      "add integer global identifier",
			key:       "age",
			val:       25,
			expectKey: "age",
			expectVal: &Integer{Val: 25},
		},
		{
			name:      "add negative integer global identifier",
			key:       "score",
			val:       -10,
			expectKey: "score",
			expectVal: &Integer{Val: -10},
		},
		{
			name:      "add float global identifier",
			key:       "height",
			val:       5.7,
			expectKey: "height",
			expectVal: &Float{Val: 5.7},
		},
		{
			name:      "add negative float global identifier",
			key:       "temperature",
			val:       -2.5,
			expectKey: "temperature",
			expectVal: &Float{Val: -2.5},
		},
		{
			name:      "add boolean true global identifier",
			key:       "isActive",
			val:       true,
			expectKey: "isActive",
			expectVal: &Boolean{Val: true},
		},
		{
			name:      "add boolean false global identifier",
			key:       "isComplete",
			val:       false,
			expectKey: "isComplete",
			expectVal: &Boolean{Val: false},
		},
		{
			name:      "add empty string global identifier",
			key:       "empty",
			val:       "",
			expectKey: "empty",
			expectVal: &String{Val: ""},
		},
		{
			name:      "add zero integer global identifier",
			key:       "zero",
			val:       0,
			expectKey: "zero",
			expectVal: &Integer{Val: 0},
		},
		{
			name:      "add zero float global identifier",
			key:       "zeroFloat",
			val:       0.0,
			expectKey: "zeroFloat",
			expectVal: &Float{Val: 0.0},
		},
		{
			name:      "add string with special characters",
			key:       "message",
			val:       "Hello, World! @#$%",
			expectKey: "message",
			expectVal: &String{Val: "Hello, World! @#$%"},
		},
		{
			name:      "add large integer global identifier",
			key:       "bigNumber",
			val:       999999999,
			expectKey: "bigNumber",
			expectVal: &Integer{Val: 999999999},
		},
		{
			name:      "add integer slice global identifier",
			key:       "numbers",
			val:       []int{1, 2, 3},
			expectKey: "numbers",
			expectVal: &Array{Elements: []Object{
				&Integer{Val: 1},
				&Integer{Val: 2},
				&Integer{Val: 3},
			}},
		},
		{
			name:      "add string slice global identifier",
			key:       "names",
			val:       []string{"Alice", "Bob", "Charlie"},
			expectKey: "names",
			expectVal: &Array{Elements: []Object{
				&String{Val: "Alice"},
				&String{Val: "Bob"},
				&String{Val: "Charlie"},
			}},
		},
		{
			name:      "add mixed type slice global identifier",
			key:       "mixed",
			val:       []any{"hello", 42, true},
			expectKey: "mixed",
			expectVal: &Array{Elements: []Object{
				&String{Val: "hello"},
				&Integer{Val: 42},
				&Boolean{Val: true},
			}},
		},
		{
			name:      "add empty slice global identifier",
			key:       "emptySlice",
			val:       []int{},
			expectKey: "emptySlice",
			expectVal: &Array{Elements: []Object{}},
		},
		{
			name:      "add object/map global identifier",
			key:       "user",
			val:       map[string]any{"name": "John", "age": 30},
			expectKey: "user",
			expectVal: NewObj(map[string]Object{
				"name": &String{Val: "John"},
				"age":  &Integer{Val: 30},
			}),
		},
		{
			name: "add nested object global identifier",
			key:  "config",
			val: map[string]any{
				"database": map[string]any{
					"host": "localhost",
					"port": 5432,
				},
			},
			expectKey: "config",
			expectVal: NewObj(map[string]Object{
				"database": NewObj(map[string]Object{
					"host": &String{Val: "localhost"},
					"port": &Integer{Val: 5432},
				}),
			}),
		},
		{
			name:      "add object with slice global identifier",
			key:       "data",
			val:       map[string]any{"tags": []string{"go", "test", "unit"}, "count": 3},
			expectKey: "data",
			expectVal: NewObj(map[string]Object{
				"tags": &Array{Elements: []Object{
					&String{Val: "go"},
					&String{Val: "test"},
					&String{Val: "unit"},
				}},
				"count": &Integer{Val: 3},
			}),
		},
		{
			name: "add slice of objects global identifier",
			key:  "users",
			val: []map[string]any{
				{"name": "Alice", "age": 25},
				{"name": "Bob", "age": 30},
			},
			expectKey: "users",
			expectVal: &Array{Elements: []Object{
				NewObj(map[string]Object{
					"name": &String{Val: "Alice"},
					"age":  &Integer{Val: 25},
				}),
				NewObj(map[string]Object{
					"name": &String{Val: "Bob"},
					"age":  &Integer{Val: 30},
				}),
			}},
		},
		{
			name:      "add nil slice global identifier",
			key:       "nilSlice",
			val:       []any(nil),
			expectKey: "nilSlice",
			expectVal: &Array{Elements: []Object{}},
		},
		{
			name:      "add empty object global identifier",
			key:       "emptyObj",
			val:       map[string]any{},
			expectKey: "emptyObj",
			expectVal: NewObj(nil),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			scope := NewScope()
			scope.AddGlobal(tc.key, tc.val)

			global, ok := scope.vars["global"]
			if !ok {
				t.Fatalf("The 'global' object not found in the scope")
			}

			obj, ok := global.(*Map)
			if !ok {
				t.Fatalf("The 'global' object is not of type Obj")
			}

			val, ok := obj.Pairs[tc.key]
			if !ok {
				t.Fatalf("The 'global' object does not have key %s", tc.key)
			}

			if val.Type() != tc.expectVal.Type() {
				t.Fatalf(
					"Expected 'global[%s]' type to be %q, got %q",
					tc.key,
					tc.expectVal.Type(),
					val.Type(),
				)
			}

			if val.String() != tc.expectVal.String() {
				t.Fatalf("Expected 'global[%s]' to be %q, got %q", tc.key, tc.expectVal, val)
			}
		})
	}
}
