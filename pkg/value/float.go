package value

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/textwire/textwire/v3/pkg/utils"
)

type Float struct {
	Val float64
}

func (f *Float) Type() ValueType {
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

func (f *Float) SubtractFromFloat(num uint) error {
	// Convert the float to a string
	strVal := strconv.FormatFloat(f.Val, 'f', -1, 64)

	if !strings.Contains(strVal, ".") {
		f.Val -= float64(num)
		return nil
	}

	nums := strings.Split(strVal, ".")

	// Parse the integer part
	intPart, err := strconv.ParseUint(nums[0], 10, 64)
	if err != nil {
		return err
	}

	// Subtract the uint value from the integer part
	intPart -= uint64(num)

	// Combine the modified integer part and the decimal part
	resultStr := fmt.Sprintf("%d.%s", intPart, nums[1])

	// Parse the result back to float64
	result, err := strconv.ParseFloat(resultStr, 64)
	if err != nil {
		return err
	}

	f.Val = result

	return nil
}

func (f *Float) Is(t ValueType) bool {
	return t == f.Type()
}
