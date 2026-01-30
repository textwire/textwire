package utils

import (
	"strconv"
	"strings"
)

// FloatToStr converts float64 to string using a
// precision of -1 to preserve the exact value
func FloatToStr(f float64) string {
	// Use 'floatStr' format first to see if it uses scientific notation
	floatStr := strconv.FormatFloat(f, 'g', -1, 64)

	// If 'floatStr' format uses scientific notation and the number is very large, use it
	if (strings.Contains(floatStr, "e") || strings.Contains(floatStr, "E")) &&
		(f > 1e20 || f < -1e20) {
		return floatStr
	}

	// Otherwise use 'f' format for decimal notation
	str := strconv.FormatFloat(f, 'f', -1, 64)

	// If no decimal point, add .0
	if !strings.Contains(str, ".") {
		str += ".0"
	}

	return str
}
