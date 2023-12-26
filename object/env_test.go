package object

import "testing"

func TestEnvFromMap(t *testing.T) {
	t.Run("with string values", func(t *testing.T) {
		m := map[string]interface{}{
			"foo": "bar",
			"baz": "qux",
		}

		env, err := EnvFromMap(m)

		if err != nil {
			t.Fatalf("EnvFromMap returned an error: %s", err)
		}

		for key, val := range m {
			obj, ok := env.Get(key)

			if !ok {
				t.Fatalf("Env.Get returned false for key %s", key)
			}

			str, ok := obj.(*String)

			if !ok {
				t.Fatalf("Env.Get returned a non-String object for key %s", key)
			}

			if str.Value != val.(string) {
				t.Fatalf("Env.Get returned a String object with an incorrect value for key %s", key)
			}
		}
	})

	t.Run("with int values", func(t *testing.T) {
		m := map[string]interface{}{
			"foo": 1,
			"baz": 2,
		}

		env, err := EnvFromMap(m)

		if err != nil {
			t.Fatalf("EnvFromMap returned an error: %s", err)
		}

		for key, val := range m {
			obj, ok := env.Get(key)

			if !ok {
				t.Fatalf("Env.Get returned false for key %s", key)
			}

			integer, ok := obj.(*Integer)

			if !ok {
				t.Fatalf("Env.Get returned a non-Integer object for key %s", key)
			}

			if integer.Value != int64(val.(int)) {
				t.Fatalf("Env.Get returned an Integer object with an incorrect value for key %s", key)
			}
		}
	})
}
