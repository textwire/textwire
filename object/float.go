package object

import (
	"fmt"
	"strconv"
	"strings"
)

type Float struct {
	Value float64
}

func (f *Float) Type() ObjectType {
	return FLOAT_OBJ
}

func (f *Float) String() string {
	return fmt.Sprintf("%g", f.Value)
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
