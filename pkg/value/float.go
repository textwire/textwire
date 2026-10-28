package value

import (
	"fmt"
	"math"

	"github.com/textwire/textwire/v5/pkg/utils"
)

type Float struct {
	Val float64
}

func (*Float) Type() ValueType {
	return FLOAT_VAL
}

func (f *Float) String() string {
	return utils.FloatToStr(f.Val)
}

func (f *Float) Dump(ident int) string {
	return fmt.Sprintf(`<span style="%s">%s</span>`, DUMP_NUM, f)
}

func (e *Float) JSON() (string, error) {
	if math.IsNaN(e.Val) || math.IsInf(e.Val, 0) {
		return "null", nil
	}
	return e.String(), nil
}

func (f *Float) Native() any {
	return f.Val
}

func (f *Float) Is(t ValueType) bool {
	return t == f.Type()
}
