package value

import (
	"bytes"
	"fmt"
)

var outputHTML = `
<div style="
  overflow-x: auto !important;
  overflow-y: hidden !important;
  scrollbar-width: thin !important;
  margin: 4px !important;
  font-size: 0.9rem !important;
">
  <pre style="
    background-color: #212121 !important;
    color: white !important;
    padding: 20px !important;
    border-radius: 5px !important;
    margin: 0 !important !important;
    width: fit-content !important;
  ">%s</pre>
</div>`

type Dump struct {
	Vals []Literal
}

func NewDump(cap int) *Dump {
	return &Dump{Vals: make([]Literal, 0, cap)}
}

func (*Dump) Type() ValueType {
	return DUMP_VAL
}

func (d *Dump) String() string {
	var out bytes.Buffer
	out.Grow(len(d.Vals))

	for _, v := range d.Vals {
		fmt.Fprintf(&out, outputHTML, v.Dump(0))
	}

	return out.String()
}

func (d *Dump) Is(t ValueType) bool {
	return t == d.Type()
}
