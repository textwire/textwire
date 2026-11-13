package value

import (
	"fmt"
	"testing"
)

func TestDump_String(t *testing.T) {
	tests := []struct {
		name string
		vals []Literal
		want string
	}{
		{
			name: "no values",
			vals: nil,
			want: "",
		},
		{
			name: "single value",
			vals: []Literal{&Str{Val: "hello"}},
			want: fmt.Sprintf(outputHTML, (&Str{Val: "hello"}).Dump(0)),
		},
		{
			name: "multiple values",
			vals: []Literal{&Str{Val: "hello"}, &Int{Val: 1}},
			want: fmt.Sprintf(outputHTML, (&Str{Val: "hello"}).Dump(0)) +
				fmt.Sprintf(outputHTML, (&Int{Val: 1}).Dump(0)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dump := &Dump{Vals: tt.vals}

			if got := dump.String(); got != tt.want {
				t.Fatalf("expect %q, got %q", tt.want, got)
			}
		})
	}
}
