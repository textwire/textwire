package object

import "bytes"

var DumpHTML = `<style>
.textwire-dump {
	overflow-x: auto;
	overflow-y: hidden;
    scrollbar-width: thin;
	margin: 4px;
}

.textwire-prop { color: #f8f8f2 }
.textwire-str { color: #c3e88d }
.textwire-num { color: #76a8ff }
.textwire-key { color: #c792ea }
.textwire-brace { color: #e99f33 }

.textwire-dump pre {
	background-color: #212121;
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
	return "dump stmt"
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
