package hashring

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"strconv"
	"sync"
)

type Server struct {
	Addr string
	Weight int
}

const (
	DefaultVirtualNodeNum = 100
)

type VirtualNode struct {
	NodeKey string
	SpotValue uint32
}

type NodeArray []VirtualNode

func (s NodeArray) Len() int {return len(s)}
func (s NodeArray) Less(i,j int) bool {return s[i].SpotValue < s[j].SpotValue}
func (s NodeArray) Swap(i,j int) {s[i].SpotValue, s[j].SpotValue = s[j].SpotValue, s[i].SpotValue}
func (s NodeArray) Sort() {sort.Sort(s)}

type HashRing struct {
	Mu *sync.RWMutex
	Servers map[string]int
	Nodes NodeArray
}

func NewHashRing(servers map[string]int) *HashRing {
	return &HashRing{
		Mu:      new(sync.RWMutex),
		Servers: servers,
		Nodes:   nil,
	}
}

func (s *HashRing) Add(server Server) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.Servers[server.Addr] = server.Weight
	s.Generate()
}

func (s *HashRing) Delete(server Server) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	delete(s.Servers, server.Addr)
	s.Generate()
}

func (s *HashRing) Generate() {
	if len(s.Servers) < 1 {
		return
	}

	//var SpotNums int
	for k,v := range s.Servers {
		totalSpotNums := DefaultVirtualNodeNum * v
		for i := 0; i < totalSpotNums; i++ {
			iString := strconv.Itoa(i)
			virtuallAddr := fmt.Sprintf("%v:%v", k, iString)

			virtualNode := &VirtualNode{
				NodeKey:   k,
				SpotValue: HashSha256(virtuallAddr),
			}
			s.Nodes = append(s.Nodes, *virtualNode)
		}
	}
}

func (s *HashRing) GetServer(w uint32) (addr string) {
	i := sort.Search(s.Nodes.Len(), func(i int) bool {
		return s.Nodes[i].SpotValue >= w
	})
	return s.Nodes[i].NodeKey
}

func HashSha256(s string) (v uint32) {
	hash := sha256.New()
	hash.Write([]byte(s))
	sum := hash.Sum(nil)

	if len(sum[6:10]) == 4 {
		v = uint32(sum[3]) << 24 | uint32(sum[2]) << 16 | uint32(sum[1]) << 8 | uint32(sum[0])
	}
	return
}







