package handling

import (
	rbt "github.com/emirpasic/gods/trees/redblacktree"
	"github.com/serialx/hashring"
	"github.com/sirupsen/logrus"
)

//    --------------------------------------------------------------------------------
//                                    STRUCTS
//    --------------------------------------------------------------------------------

type HashRing struct {
	Ring  *hashring.HashRing
	Trees map[string]*rbt.Tree // Is map[nodeAddr]*RedBlackTree
	//Addresses []string             // Addresses of working nodes
}

//    --------------------------------------------------------------------------------
//                                     METHODS
//    --------------------------------------------------------------------------------

// hashRingInit returns an empty structure *HashRing
func hashRingInit() *HashRing {
	var memcacheServers []string
	ring := hashring.New(memcacheServers)
	trees := make(map[string]*rbt.Tree)

	return &HashRing{
		Ring:  ring,
		Trees: trees,
	}
}

// calculateNodeAddr calculates the address of the service based on the unique message key
func (h *HashRing) calculateNodeAddr(msgKey string) (server string, ok bool) {
	server, ok = h.Ring.GetNode(msgKey)

	logrus.Infof("Вычисляем адрес ноды с помощью shring. Server = %v , ok = %v;;; size = %v\n", server, ok, h.Ring.Size())

	return server, ok
}

// addNodeToHashRing adds the address of the new service to the existing ones
func (h *HashRing) addNodeToHashRing(nodeAddr string) {
	newRing := h.Ring.AddNode(nodeAddr)
	h.Ring = newRing
}

// removeNodeFromHashRing removes the service address from the ring
func (h *HashRing) removeNodeFromHashRing(nodeAddr string) {
	newRing := h.Ring.RemoveNode(nodeAddr)
	h.Ring = newRing
}

//	memcacheServers := []string{
//		"191.0.0.0:3000",
//		"192.0.0.0:3000",
//		"193.0.0.0:3000",
//	}
//
//	ring := hashring.New(memcacheServers)
//
//	//ring = ring.RemoveNode("193.0.0.0:3000")
//
//	server, _ := ring.GetNode("NEWKEY1")
//
//	//ring = ring.AddNode("193.0.0.0:3000")
//
//	fmt.Printf("For key  %v node addr1 %v \n\n", "NEWKEY1", server)
//	fmt.Println("sdf")
