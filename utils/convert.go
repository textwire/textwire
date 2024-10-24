package utils

import "strconv"

// FloatToStr converts float64 to string using a
// precision of -1 to preserve the exact value
func FloatToStr(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}
