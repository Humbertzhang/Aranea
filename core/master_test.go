package core

import "testing"

func TestMaster_PingOnce(t *testing.T) {
	master := &Master{}
	master.pingOnce("https://www.baidu.com")
}
