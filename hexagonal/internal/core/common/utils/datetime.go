package utils

import "time"

func GetUTCCurrentMillis() uint64 {
	return uint64(time.Now().UTC().UnixMilli())
}
