package value

import (
	"reflect"
	"testing"
)

func assertValue(t *testing.T, inp any, expect Literal) {
	t.Helper()

	got := NativeToValue(inp)
	if !reflect.DeepEqual(got, expect) {
		t.Fatalf("NativeToValue(%T)\n got: %#v\nexpect: %#v", inp, got, expect)
	}
}

func assertErrorMessage(t *testing.T, inp any, expect string) {
	t.Helper()

	got := NativeToValue(inp)

	errObj, ok := got.(*Error)
	if !ok {
		t.Fatalf("NativeToValue(%T) returned %T, expect *Error", inp, got)
	}

	if msg := errObj.Err.Message(); msg != expect {
		t.Fatalf("error message\n got: %q\nexpect: %q", msg, expect)
	}
}
