package value

import (
	"bytes"
	"fmt"
	"strings"
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
	Vals []string
}

func (d *Dump) Type() ValueType {
	return DUMP_OBJ
}

func (d *Dump) String() string {
	var out bytes.Buffer
	for _, v := range d.Vals {
		fmt.Fprintf(&out, outputHTML, v)
	}

	return out.String()
}

func (d *Dump) Dump(ident int) string {
	return fmt.Sprintf("@dump(%s)", strings.Join(d.Vals, ", "))
}

func (d *Dump) JSON() (string, error) {
	return "", nil
}

func (d *Dump) Native() any {
	var vals []any
	for i := range d.Vals {
		vals = append(vals, d.Vals[i])
	}
	return vals
}

func (d *Dump) Is(t ValueType) bool {
	return t == d.Type()
}
