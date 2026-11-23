package value

import (
	"testing"
)

func TestNewScopeFromMap(t *testing.T) {
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
		"weight":   5.731,
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
		"rates32":  []float32{5.731, 5.731, 5.731},
	}

	expect := map[string]Value{
		"title":    &Str{Val: "Hello, World!"},
		"n":        &Int{Val: -1},
		"num":      &Int{Val: -2},
		"num1":     &Int{Val: -3},
		"num2":     &Int{Val: -4},
		"num3":     &Int{Val: -5},
		"un":       &Int{Val: 1},
		"unum":     &Int{Val: 2},
		"unum1":    &Int{Val: 3},
		"unum2":    &Int{Val: 4},
		"unum3":    &Int{Val: 5},
		"bool1":    &Bool{Val: true},
		"bool2":    &Bool{Val: false},
		"height":   &Float{Val: 5.7},
		"weight":   &Float{Val: 5.731},
		"nothing":  new(Nil),
		"ages":     &Arr{Elements: []Literal{&Int{Val: 1}, &Int{Val: 2}, &Int{Val: 3}}},
		"ages64":   &Arr{Elements: []Literal{&Int{Val: 1}, &Int{Val: 2}, &Int{Val: 3}}},
		"ages32":   &Arr{Elements: []Literal{&Int{Val: 1}, &Int{Val: 2}, &Int{Val: 3}}},
		"ages16":   &Arr{Elements: []Literal{&Int{Val: 1}, &Int{Val: 2}, &Int{Val: 3}}},
		"ages8":    &Arr{Elements: []Literal{&Int{Val: 1}, &Int{Val: 2}, &Int{Val: 3}}},
		"nums":     &Arr{Elements: []Literal{&Int{Val: 1}, &Int{Val: 2}, &Int{Val: 3}}},
		"nums64":   &Arr{Elements: []Literal{&Int{Val: 1}, &Int{Val: 2}, &Int{Val: 3}}},
		"nums32":   &Arr{Elements: []Literal{&Int{Val: 1}, &Int{Val: 2}, &Int{Val: 3}}},
		"nums16":   &Arr{Elements: []Literal{&Int{Val: 1}, &Int{Val: 2}, &Int{Val: 3}}},
		"nums8":    &Arr{Elements: []Literal{&Int{Val: 1}, &Int{Val: 2}, &Int{Val: 3}}},
		"names":    &Arr{Elements: []Literal{&Str{Val: "John"}, &Str{Val: "Jane"}}},
		"statuses": &Arr{Elements: []Literal{&Bool{Val: false}, &Bool{Val: true}}},
		"rates64": &Arr{
			Elements: []Literal{&Float{Val: 23.4}, &Float{Val: 56.7}, &Float{Val: 89.0}},
		},
		"values": &Arr{
			Elements: []Literal{&Float{Val: 23.4}, &Float{Val: 56.7}, &Float{Val: 89.0}},
		},
		"rates32": &Arr{
			Elements: []Literal{
				&Float{Val: 5.731},
				&Float{Val: 5.731},
				&Float{Val: 5.731},
			},
		},
	}

	scope, err := NewScopeFromMap(data, nil)
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

func TestNewChildScope(t *testing.T) {
	parent := NewScope(nil)
	parent.vars["title"] = &Str{Val: "Hello"}

	child := NewChildScope(parent)

	obj, ok := child.Get("title")
	if !ok {
		t.Fatalf("child.Get(title) returned !ok, expected parent's variable")
	}

	if obj.String() != "Hello" {
		t.Fatalf("child.Get(title) returned %q, expected %q", obj, "Hello")
	}
}

func TestScopeOf(t *testing.T) {
	grandparent := NewScope(nil)
	grandparent.vars["title"] = &Str{Val: "grandparent"}
	grandparent.vars["name"] = &Str{Val: "grandparent"}

	parent := NewChildScope(grandparent)
	parent.vars["name"] = &Str{Val: "parent"}

	child := NewChildScope(parent)
	child.vars["age"] = &Int{Val: 1}

	cases := []struct {
		name   string
		scope  *Scope
		key    string
		expect *Scope
	}{
		{
			name:   "key defined in current scope",
			scope:  child,
			key:    "age",
			expect: child,
		},
		{
			name:   "key defined only in grandparent scope",
			scope:  child,
			key:    "title",
			expect: grandparent,
		},
		{
			name:   "closest scope wins when key is shadowed",
			scope:  child,
			key:    "name",
			expect: parent,
		},
		{
			name:   "key not defined in any scope",
			scope:  child,
			key:    "missing",
			expect: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.scope.scopeOf(tc.key)
			if got != tc.expect {
				t.Fatalf("scopeOf(%q) returned %p, expected %p", tc.key, got, tc.expect)
			}
		})
	}
}
