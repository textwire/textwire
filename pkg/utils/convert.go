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

func ToCamel(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}

	var n strings.Builder
	n.Grow(len(s))

	capNext := false
	prevIsCap := false

	for i, v := range []byte(s) {
		vIsCap := v >= 'A' && v <= 'Z'
		vIsLow := v >= 'a' && v <= 'z'

		if capNext {
			if vIsLow {
				v += 'A'
				v -= 'a'
			}
		} else if i == 0 {
			if vIsCap {
				v += 'a'
				v -= 'A'
			}
		} else if prevIsCap && vIsCap {
			v += 'a'
			v -= 'A'
		}

		prevIsCap = vIsCap

		if vIsCap || vIsLow {
			n.WriteByte(v)
			capNext = false
		} else if vIsNum := v >= '0' && v <= '9'; vIsNum {
			n.WriteByte(v)
			capNext = true
		} else {
			capNext = v == '_' || v == ' ' || v == '-' || v == '.'
		}
	}
	return n.String()
}
