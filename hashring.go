package hashring

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"strconv"
	"sync"
)

const (
	mod = 1 >> 32 -1
	DefaultNodeNum = 100
)

type Server struct {
	Addr string
	Weight int
}

//当前服务器相关信息
var Servers []Server

func HashSha256(key string) (value uint32) {
	hash := sha256.New()
	defer hash.Reset()

	hash.Write([]byte(key))
	hashBytes  := hash.Sum(nil)

	if len(hashBytes[6:10]) == 4 {
		value = uint32(hashBytes[3]) << 24 | uint32(hashBytes[2]) << 16 | uint32(hashBytes[1]) << 8 | uint32(hashBytes[0])
	}
	return
}

type VirtualNode struct {
	nodeKey string
	spotValue uint32
}

type Nodes struct {
	VirtualNodesArray []VirtualNode
}

func (p *Nodes) Len() int {return len(p.VirtualNodesArray)}
func (p *Nodes) Less(i,j int) bool {return p.VirtualNodesArray[i].spotValue < p.VirtualNodesArray[j].spotValue}
func (p *Nodes) Swap(i,j int) {p.VirtualNodesArray[i], p.VirtualNodesArray[j] = p.VirtualNodesArray[j], p.VirtualNodesArray[i]}
func (p *Nodes) Sort() {sort.Sort(p)}


func (p *Nodes) SetVirtualNodesArray(servers []Server) {
	if len(servers) < 1 {
		return
	}

	//根据权重与节点数，维护一个map - 所有的hash圈上的值对应ip
	for _, v := range servers {
		//第一步计算出每台机器对应的虚拟节点数
		totalVirtualNodeNum := DefaultNodeNum * v.Weight
		for i := 0; i < totalVirtualNodeNum; i++ {
			iString := strconv.Itoa(i)
			//虚拟节点地址
			virtualAddr := fmt.Sprintf("%s:%s", v.Addr, iString)

			virNode := VirtualNode{
				nodeKey: v.Addr,
				spotValue: HashSha256(virtualAddr),
			}

			p.VirtualNodesArray = append(p.VirtualNodesArray, virNode)
		}
		p.Sort()
	}

}

//获取当前数据key对应的存储服务器
func (p *Nodes) getNodeSever(w uint32) (addr string) {
	i := sort.Search(len(p.VirtualNodesArray), func(i int) bool { return p.VirtualNodesArray[i].spotValue >= w })
	return p.VirtualNodesArray[i].nodeKey
}


type HashRing struct {
	Ring map[uint32]VirtualNode
	Mu *sync.RWMutex
}

func NewHashRing(spotNum int) *HashRing {
	if spotNum == 0 {
		spotNum = DefaultNodeNum
	}
	return &HashRing{
		Ring: make(map[uint32]VirtualNode, spotNum),
		Mu:   new(sync.RWMutex),
	}
}

func (s *HashRing) AddNode(node *VirtualNode) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.Ring[node.spotValue] = *node
}

func (s *HashRing) DeleteNode(node *VirtualNode) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	delete(s.Ring, node.spotValue)
}

