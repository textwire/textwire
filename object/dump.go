package object

import "bytes"

var DumpHTML = `<style>
.textwire-dump {
	overflow-x: auto;
	overflow-y: hidden;
    scrollbar-width: thin;
}

.textwire-dump pre {
	background-color: #0b0019;
	color: white;
	padding: 13px;
	border-radius: 5px;
	margin: 0 !important;
	width: fit-content;
}
</style>
<div class="textwire-dump"><pre>%s</pre></div>`

type Dump struct {
	Values []string
}

func (d *Dump) Type() ObjectType {
	return DUMP_OBJ
}

func (d *Dump) String() string {
	var out bytes.Buffer

	for _, v := range d.Values {
		out.WriteString(v)
	}

	return out.String()
}

func (d *Dump) Dump(ident int) string {
	return "dump"
}

func (d *Dump) Val() interface{} {
	var result []interface{}

	for _, v := range d.Values {
		result = append(result, v)
	}

	return result
}

func (d *Dump) Is(t ObjectType) bool {
	return t == d.Type()
}
