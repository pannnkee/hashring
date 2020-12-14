package hashring

import (
	"fmt"
	"testing"
)

func TestHashRing(t *testing.T) {

	server := make(map[string]int)
	//server["127.0.0.1"] = 1
	//server["127.0.0.2"] = 2
	//server["127.0.0.3"] = 3

	ring := NewHashRing(server)
	ring.Add(Server{
		Addr:   "127.0.0.1",
		Weight: 1,
	})
	ring.Add(Server{
		Addr:   "127.0.0.2",
		Weight: 2,
	})
	ring.Add(Server{
		Addr:   "127.0.0.3",
		Weight: 3,
	})
	sha256 := HashSha256("1.jpg")
	getServer := ring.GetServer(sha256)
	fmt.Println(getServer)
}
