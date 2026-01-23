package object

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/textwire/textwire/v2/utils"
)

type Float struct {
	Value float64
}

func (f *Float) Type() ObjectType {
	return FLOAT_OBJ
}

func (f *Float) String() string {
	return utils.FloatToStr(f.Value)
}

func (f *Float) Dump(ident int) string {
	return "<span class='textwire-num'>" + f.String() + "</span>"
}

func (f *Float) Val() any {
	return f.Value
}

func (f *Float) SubtractFromFloat(num uint) error {
	// Convert the float to a string
	strValue := strconv.FormatFloat(f.Value, 'f', -1, 64)

	if !strings.Contains(strValue, ".") {
		f.Value -= float64(num)
		return nil
	}

	nums := strings.Split(strValue, ".")

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

	f.Value = result

	return nil
}

func (f *Float) Is(t ObjectType) bool {
	return t == f.Type()
}
