package utils

import (
	"fmt"
	"testing"
)

func TestStringTimeStampNanoSecond(t *testing.T) {
	for i:= 0; i < 100; i++ {
		s := StringTimeStampNanoSecond()
		fmt.Println(s)
	}
}
