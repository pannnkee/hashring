package hashring

import (
	"fmt"
	"testing"
)

func TestHashRing(t *testing.T) {
	vNodes := new(Nodes)

	var servers []Server
	servers = append(servers, Server{"127.0.0.1", 1}, Server{"127.0.0.2", 2}, Server{"127.0.0.3", 3})

	vNodes.SetVirtualNodesArray(servers)

	fname := "2.jpg"
	sha256 := HashSha256(fname)
	ip := vNodes.getNodeSever(sha256)
	fmt.Println(ip)
}
