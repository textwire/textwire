package object

import "testing"

func TestEnvFromMap(t *testing.T) {
	vars := map[string]interface{}{
		"title":  "Hello, World!",
		"n":      -1,
		"num":    int8(-2),
		"num1":   int16(-3),
		"num2":   int32(-4),
		"num3":   int64(-5),
		"un":     uint(1),
		"unum":   uint8(2),
		"unum1":  uint16(3),
		"unum2":  uint32(4),
		"unum3":  uint64(5),
		"bool1":  true,
		"bool2":  false,
		"height": 5.7,
		"weight": float32(62.2),
	}

	expect := map[string]Object{
		"title":  &String{Value: "Hello, World!"},
		"n":      &Int{Value: -1},
		"num":    &Int8{Value: -2},
		"num1":   &Int16{Value: -3},
		"num2":   &Int32{Value: -4},
		"num3":   &Int64{Value: -5},
		"un":     &Uint{Value: 1},
		"unum":   &Uint8{Value: 2},
		"unum1":  &Uint16{Value: 3},
		"unum2":  &Uint32{Value: 4},
		"unum3":  &Uint64{Value: 5},
		"bool1":  &Boolean{Value: true},
		"bool2":  &Boolean{Value: false},
		"height": &Float64{Value: 5.7},
		"weight": &Float32{Value: 62.2},
	}

	env, err := EnvFromMap(vars)

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
