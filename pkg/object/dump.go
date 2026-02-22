package object

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
	Values []string
}

func (d *Dump) Type() ObjectType {
	return DUMP_OBJ
}

func (d *Dump) String() string {
	var out bytes.Buffer
	for _, v := range d.Values {
		fmt.Fprintf(&out, outputHTML, v)
	}

	return out.String()
}

func (d *Dump) Dump(ident int) string {
	return fmt.Sprintf("@dump(%s)", strings.Join(d.Values, ", "))
}

func (d *Dump) Val() any {
	var values []any
	for i := range d.Values {
		values = append(values, d.Values[i])
	}
	return values
}

func (d *Dump) Is(t ObjectType) bool {
	return t == d.Type()
}
