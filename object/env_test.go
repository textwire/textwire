package object

import (
	"testing"
)

func TestEnvFromMap(t *testing.T) {
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
		"title":    &Str{Value: "Hello, World!"},
		"n":        &Int{Value: -1},
		"num":      &Int{Value: -2},
		"num1":     &Int{Value: -3},
		"num2":     &Int{Value: -4},
		"num3":     &Int{Value: -5},
		"un":       &Int{Value: 1},
		"unum":     &Int{Value: 2},
		"unum1":    &Int{Value: 3},
		"unum2":    &Int{Value: 4},
		"unum3":    &Int{Value: 5},
		"bool1":    &Bool{Value: true},
		"bool2":    &Bool{Value: false},
		"height":   &Float{Value: 5.7},
		"weight":   &Float{Value: float64(float32val)},
		"nothing":  &Nil{},
		"ages":     &Array{Elements: []Object{&Int{Value: 1}, &Int{Value: 2}, &Int{Value: 3}}},
		"ages64":   &Array{Elements: []Object{&Int{Value: 1}, &Int{Value: 2}, &Int{Value: 3}}},
		"ages32":   &Array{Elements: []Object{&Int{Value: 1}, &Int{Value: 2}, &Int{Value: 3}}},
		"ages16":   &Array{Elements: []Object{&Int{Value: 1}, &Int{Value: 2}, &Int{Value: 3}}},
		"ages8":    &Array{Elements: []Object{&Int{Value: 1}, &Int{Value: 2}, &Int{Value: 3}}},
		"nums":     &Array{Elements: []Object{&Int{Value: 1}, &Int{Value: 2}, &Int{Value: 3}}},
		"nums64":   &Array{Elements: []Object{&Int{Value: 1}, &Int{Value: 2}, &Int{Value: 3}}},
		"nums32":   &Array{Elements: []Object{&Int{Value: 1}, &Int{Value: 2}, &Int{Value: 3}}},
		"nums16":   &Array{Elements: []Object{&Int{Value: 1}, &Int{Value: 2}, &Int{Value: 3}}},
		"nums8":    &Array{Elements: []Object{&Int{Value: 1}, &Int{Value: 2}, &Int{Value: 3}}},
		"names":    &Array{Elements: []Object{&Str{Value: "John"}, &Str{Value: "Jane"}}},
		"statuses": &Array{Elements: []Object{&Bool{Value: false}, &Bool{Value: true}}},
		"rates64":  &Array{Elements: []Object{&Float{Value: 23.4}, &Float{Value: 56.7}, &Float{Value: 89.0}}},
		"values":   &Array{Elements: []Object{&Float{Value: 23.4}, &Float{Value: 56.7}, &Float{Value: 89.0}}},
		"rates32":  &Array{Elements: []Object{&Float{Value: float64(float32val)}, &Float{Value: float64(float32val)}, &Float{Value: float64(float32val)}}},
	}

	env, err := EnvFromMap(data)
	if err != nil {
		t.Fatalf("EnvFromMap returned an error: %s", err)
	}

	for key, val := range expect {
		obj, ok := env.Get(key)

		if !ok {
			t.Fatalf("Env.Get(%s) returned !ok", key)
		}

		if obj.Type() != val.Type() {
			t.Fatalf("Env.Get(%s) returned %s, expect %s", key, obj.Type(), val.Type())
		}

		if obj.String() != val.String() {
			t.Fatalf("Env.Get(%s) returned %s, expect %s", key, obj.String(), val.String())
		}
	}
}

func TestAddGlobalVar(t *testing.T) {
	cases := []struct {
		name      string
		key       string
		val       any
		expectKey string
		expectVal Object
	}{
		{
			name:      "add string global variable",
			key:       "name",
			val:       "Amy Adams",
			expectKey: "name",
			expectVal: &Str{Value: "Amy Adams"},
		},
		{
			name:      "add integer global variable",
			key:       "age",
			val:       25,
			expectKey: "age",
			expectVal: &Int{Value: 25},
		},
		{
			name:      "add negative integer global variable",
			key:       "score",
			val:       -10,
			expectKey: "score",
			expectVal: &Int{Value: -10},
		},
		{
			name:      "add float global variable",
			key:       "height",
			val:       5.7,
			expectKey: "height",
			expectVal: &Float{Value: 5.7},
		},
		{
			name:      "add negative float global variable",
			key:       "temperature",
			val:       -2.5,
			expectKey: "temperature",
			expectVal: &Float{Value: -2.5},
		},
		{
			name:      "add boolean true global variable",
			key:       "isActive",
			val:       true,
			expectKey: "isActive",
			expectVal: &Bool{Value: true},
		},
		{
			name:      "add boolean false global variable",
			key:       "isComplete",
			val:       false,
			expectKey: "isComplete",
			expectVal: &Bool{Value: false},
		},
		{
			name:      "add empty string global variable",
			key:       "empty",
			val:       "",
			expectKey: "empty",
			expectVal: &Str{Value: ""},
		},
		{
			name:      "add zero integer global variable",
			key:       "zero",
			val:       0,
			expectKey: "zero",
			expectVal: &Int{Value: 0},
		},
		{
			name:      "add zero float global variable",
			key:       "zeroFloat",
			val:       0.0,
			expectKey: "zeroFloat",
			expectVal: &Float{Value: 0.0},
		},
		{
			name:      "add string with special characters",
			key:       "message",
			val:       "Hello, World! @#$%",
			expectKey: "message",
			expectVal: &Str{Value: "Hello, World! @#$%"},
		},
		{
			name:      "add large integer global variable",
			key:       "bigNumber",
			val:       999999999,
			expectKey: "bigNumber",
			expectVal: &Int{Value: 999999999},
		},
		{
			name:      "add integer slice global variable",
			key:       "numbers",
			val:       []int{1, 2, 3},
			expectKey: "numbers",
			expectVal: &Array{Elements: []Object{
				&Int{Value: 1},
				&Int{Value: 2},
				&Int{Value: 3},
			}},
		},
		{
			name:      "add string slice global variable",
			key:       "names",
			val:       []string{"Alice", "Bob", "Charlie"},
			expectKey: "names",
			expectVal: &Array{Elements: []Object{
				&Str{Value: "Alice"},
				&Str{Value: "Bob"},
				&Str{Value: "Charlie"},
			}},
		},
		{
			name:      "add mixed type slice global variable",
			key:       "mixed",
			val:       []any{"hello", 42, true},
			expectKey: "mixed",
			expectVal: &Array{Elements: []Object{
				&Str{Value: "hello"},
				&Int{Value: 42},
				&Bool{Value: true},
			}},
		},
		{
			name:      "add empty slice global variable",
			key:       "emptySlice",
			val:       []int{},
			expectKey: "emptySlice",
			expectVal: &Array{Elements: []Object{}},
		},
		{
			name:      "add object/map global variable",
			key:       "user",
			val:       map[string]any{"name": "John", "age": 30},
			expectKey: "user",
			expectVal: NewObj(map[string]Object{
				"name": &Str{Value: "John"},
				"age":  &Int{Value: 30},
			}),
		},
		{
			name: "add nested object global variable",
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
					"host": &Str{Value: "localhost"},
					"port": &Int{Value: 5432},
				}),
			}),
		},
		{
			name:      "add object with slice global variable",
			key:       "data",
			val:       map[string]any{"tags": []string{"go", "test", "unit"}, "count": 3},
			expectKey: "data",
			expectVal: NewObj(map[string]Object{
				"tags": &Array{Elements: []Object{
					&Str{Value: "go"},
					&Str{Value: "test"},
					&Str{Value: "unit"},
				}},
				"count": &Int{Value: 3},
			}),
		},
		{
			name: "add slice of objects global variable",
			key:  "users",
			val: []map[string]any{
				{"name": "Alice", "age": 25},
				{"name": "Bob", "age": 30},
			},
			expectKey: "users",
			expectVal: &Array{Elements: []Object{
				NewObj(map[string]Object{
					"name": &Str{Value: "Alice"},
					"age":  &Int{Value: 25},
				}),
				NewObj(map[string]Object{
					"name": &Str{Value: "Bob"},
					"age":  &Int{Value: 30},
				}),
			}},
		},
		{
			name:      "add nil slice global variable",
			key:       "nilSlice",
			val:       []any(nil),
			expectKey: "nilSlice",
			expectVal: &Array{Elements: []Object{}},
		},
		{
			name:      "add empty object global variable",
			key:       "emptyObj",
			val:       map[string]any{},
			expectKey: "emptyObj",
			expectVal: NewObj(nil),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			env := NewEnv()
			env.AddGlobalVar(tc.key, tc.val)

			global, exists := env.store["global"]
			if !exists {
				t.Fatalf("'global' object not found in environment")
			}

			obj, isObj := global.(*Obj)
			if !isObj {
				t.Fatalf("'global' object is not of type Obj")
			}

			val, hasVal := obj.Pairs[tc.key]
			if !hasVal {
				t.Fatalf("'global' object does not have key %s", tc.key)
			}

			if val.Type() != tc.expectVal.Type() {
				t.Fatalf("Expected 'global[%s]' type to be %s, got %s",
					tc.key, tc.expectVal.Type(), val.Type())
			}

			if val.String() != tc.expectVal.String() {
				t.Fatalf("Expected 'global[%s]' to be '%s', got '%s'",
					tc.key, tc.expectVal, val)
			}
		})
	}
}
