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
			t.Fatalf("Env.Get(%s) returned %s, expected %s", key, obj.Type(), val.Type())
		}

		if obj.String() != val.String() {
			t.Fatalf("Env.Get(%s) returned %s, expected %s", key, obj.String(), val.String())
		}
	}
}
