package object

import (
	"bytes"
	"fmt"
	"strings"
)

var outputHTML = `<style>
.textwire-dump {
	overflow-x: auto;
	overflow-y: hidden;
    scrollbar-width: thin;
	margin: 4px;
}

.textwire-prop { color: #f8f8f2 }
.textwire-str { color: #c3e88d }
.textwire-num { color: #76a8ff }
.textwire-keyword { color: #c792ea }
.textwire-brace { color: #e99f33 }
.textwire-meta { color:#2c8ed0 }

.textwire-dump pre {
	background-color: #212121;
	color: white;
	padding: 20px;
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
		fmt.Fprintf(&out, outputHTML, v)
	}

	return out.String()
}

func (d *Dump) Dump(ident int) string {
	return fmt.Sprintf("@dump(%s)", strings.Join(d.Values, ", "))
}

func (d *Dump) Val() any {
	var res []any
	for i := range d.Values {
		res = append(res, d.Values[i])
	}
	return res
}

func (d *Dump) Is(t ObjectType) bool {
	return t == d.Type()
}
