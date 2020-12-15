package hashring

import (
	"fmt"
	"testing"
)

func TestHashRing(t *testing.T) {

	server := make(map[string]int)
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

	fmt.Println(ring.GetServer(HashSha256("1")))
	fmt.Println(ring.GetServer(HashSha256("1.jpg")))
	fmt.Println(ring.GetServer(HashSha256("666.png")))
	fmt.Println(ring.GetServer(HashSha256("pannnkee.jpg")))

	fmt.Println("delete")

	ring.Delete(Server{
		Addr:   "127.0.0.3",
		Weight: 3,
	})
	fmt.Println(ring.GetServer(HashSha256("1")))
	fmt.Println(ring.GetServer(HashSha256("1.jpg")))
	fmt.Println(ring.GetServer(HashSha256("666.png")))
	fmt.Println(ring.GetServer(HashSha256("pannnkee.jpg")))

}
