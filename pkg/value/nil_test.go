package value

import "testing"

func TestNilOrigin(t *testing.T) {
	cases := []struct {
		name   string
		nilObj *Nil
		origin NilOrigin
	}{
		{"zero value has unset origin", &Nil{}, NilOriginUnset},
		{
			"new nil keeps missing property origin",
			NewNil(NilOriginMissingProp),
			NilOriginMissingProp,
		},
		{
			"new nil keeps missing index origin",
			NewNil(NilOriginMissingIndex),
			NilOriginMissingIndex,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.nilObj.Origin(); got != tc.origin {
				t.Fatalf("expect origin %q, got %q", tc.origin, got)
			}
		})
	}
}

func TestNilStaysNil(t *testing.T) {
	nilObj := NewNil(NilOriginMissingProp)

	if got := nilObj.String(); got != "" {
		t.Errorf("expect empty string, got %q", got)
	}

	if got, err := nilObj.JSON(); err != nil || got != "null" {
		t.Errorf("expect JSON null, got %q, %v", got, err)
	}

	if got := nilObj.Native(); got != nil {
		t.Errorf("expect native nil, got %v", got)
	}

	if !nilObj.Is(NIL_VAL) {
		t.Errorf("expect nil value type")
	}
}
