package helper

import (
	"strconv"
)

// StringToInt64 :nodoc:
func StringToInt64(val string) int64 {
	num, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0
	}

	return num
}
