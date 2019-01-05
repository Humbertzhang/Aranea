package core

import "testing"

func TestMaster_PingOnce(t *testing.T) {
	master := &Master{}
	master.pingOnce("http://127.0.0.1:9699/node/pong")
}

