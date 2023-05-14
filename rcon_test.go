package mcprotocolqr

import (
	"testing"
	"time"
)

func TestRconConnect(t *testing.T) {
	rs, err := NewRconServer("127.0.0.1", "10304", "password", 5*time.Second)
	if err != nil {
		panic(err)
	}
	t.Log(rs)

	rs.Run("list")
	t.Log(rs.Get())
	t.Log("Success")

}
