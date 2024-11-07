package utils

import "strconv"

func StrIsInt(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}
