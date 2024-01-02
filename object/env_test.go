package object

import "testing"

func TestEnvFromMap(t *testing.T) {
	var double float32 = 5.7

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
		"weight": double,
	}

	expect := map[string]Object{
		"title":  &String{Value: "Hello, World!"},
		"n":      &Int{Value: -1},
		"num":    &Int{Value: -2},
		"num1":   &Int{Value: -3},
		"num2":   &Int{Value: -4},
		"num3":   &Int{Value: -5},
		"un":     &Int{Value: 1},
		"unum":   &Int{Value: 2},
		"unum1":  &Int{Value: 3},
		"unum2":  &Int{Value: 4},
		"unum3":  &Int{Value: 5},
		"bool1":  &Boolean{Value: true},
		"bool2":  &Boolean{Value: false},
		"height": &Float{Value: 5.7},
		"weight": &Float{Value: float64(double)},
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
