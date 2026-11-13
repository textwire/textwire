package value

import (
	"testing"
)

func TestRaw(t *testing.T) {
	t.Run("marks string as raw", func(t *testing.T) {
		str := &Str{Val: "<b>"}
		got := Raw(str).(*Str)

		if !got.raw || got.Val != "<b>" {
			t.Fatalf("expect raw string %q, got %#v", "<b>", got)
		}

		if str.raw {
			t.Fatal("expect the original string to stay unmarked")
		}
	})

	t.Run("returns raw string as is", func(t *testing.T) {
		str := &Str{Val: "<b>", raw: true}

		if got := Raw(str); got != str {
			t.Fatalf("expect the same string back, got %#v", got)
		}
	})

	t.Run("returns other literal types as is", func(t *testing.T) {
		values := []Literal{
			&Int{Val: 1},
			&Float{Val: 1.5},
			&Bool{Val: true},
			&Bool{Val: false},
			new(Nil),
		}

		for _, val := range values {
			if got := Raw(val); got != val {
				t.Fatalf("expect %#v back, got %#v", val, got)
			}
		}
	})

	t.Run("copies array with strings marked as raw", func(t *testing.T) {
		arr := &Arr{Elements: []Literal{&Str{Val: "<b>"}, &Int{Val: 1}}}
		got := Raw(arr).(*Arr)

		if !got.Elements[0].(*Str).raw {
			t.Fatal("expect the string in the copy to be raw")
		}

		if arr.Elements[0].(*Str).raw {
			t.Fatal("expect the string in the original array to stay unmarked")
		}
	})

	t.Run("copies nested object with strings marked as raw", func(t *testing.T) {
		inner := NewObj(map[string]Literal{"name": &Str{Val: "<b>"}})
		obj := NewObj(map[string]Literal{"user": inner})
		gotInner := Raw(obj).(*Obj).Pairs["user"].(*Obj)

		if !gotInner.Pairs["name"].(*Str).raw {
			t.Fatal("expect the nested string in the copy to be raw")
		}

		if inner.Pairs["name"].(*Str).raw {
			t.Fatal("expect the nested string in the original object to stay unmarked")
		}
	})

	t.Run("keeps shared array as one copy", func(t *testing.T) {
		shared := &Arr{Elements: []Literal{&Str{Val: "<b>"}}}
		got := Raw(&Arr{Elements: []Literal{shared, shared}}).(*Arr)

		if got.Elements[0] == shared || got.Elements[0] != got.Elements[1] {
			t.Fatal("expect both elements to be the same copy of the shared array")
		}
	})

	t.Run("copies array that contains itself", func(t *testing.T) {
		arr := &Arr{Elements: []Literal{&Str{Val: "<b>"}}}
		arr.Elements = append(arr.Elements, arr)
		got := Raw(arr).(*Arr)

		if got == arr || got.Elements[1] != got {
			t.Fatal("expect the copy to contain itself, like the original does")
		}
	})
}
