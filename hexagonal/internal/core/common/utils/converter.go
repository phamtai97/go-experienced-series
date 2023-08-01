package utils

import "strconv"

func ConvertUInt64ToString(number uint64) string {
	return strconv.FormatUint(number, 10)
}
