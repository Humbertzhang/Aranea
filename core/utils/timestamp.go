package utils

import (
	"strconv"
	"time"
)

func StringTimeStampNanoSecond() string {
	now := time.Now()
	nanos := int(now.UnixNano())
	s := strconv.Itoa(nanos)
	return s
}
