package handling

import (
	rbt "github.com/emirpasic/gods/trees/redblacktree"
	"github.com/serialx/hashring"
)

type HashingInfo struct {
	Ring      *hashring.HashRing
	Trees     map[string]*rbt.Tree // Is map[nodeAddr]*RedBlackTree
	Addresses []string             // Addresses of working nodes
}

//// computing a node for key distribution
//
//func ConsictanceInit() *Consictance {
//	ring := hashring.New()
//
//	return nil
//}

//func s() {
//	memcacheServers := []string{
//		"191.0.0.0:3000",
//		//"192.0.0.0:3000",
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
//
//}
//
//func RedBlackTree() {
//	//tree := rbt.NewWithIntComparator() // empty (keys are of type int)
//	tree := rbt.NewWithStringComparator()
//
//	tree.Put("1f", "x") // 1->x
//	tree.Put("2f", "b") // 1->x, 2->b (in order)
//	tree.Put("1f", "a") // 1->a, 2->b (in order, replacement)
//	tree.Put("3f", "c") // 1->a, 2->b, 3->c (in order)
//	tree.Put("4f", "d") // 1->a, 2->b, 3->c, 4->d (in order)
//	tree.Put("5f", "e") // 1->a, 2->b, 3->c, 4->d, 5->e (in order)
//	tree.Put("6f", "f") // 1->a, 2->b, 3->c, 4->d, 5->e, 6->f (in order)
//
//	fmt.Println(tree)
//	//
//	//  RedBlackTree
//	//  │           ┌── 6
//	//	│       ┌── 5
//	//	│   ┌── 4
//	//	│   │   └── 3
//	//	└── 2
//	//		└── 1
//
//	v, found := tree.Get("6")
//	logrus.Infof("V : %#v , FOUND : %#v \n\n", v, found)
//
//	v, found = tree.Get("6f")
//	logrus.Infof("V : %#v , FOUND : %#v \n\n", v, found)
//
//	_ = tree.Values() // []interface {}{"a", "b", "c", "d", "e", "f"} (in order)
//	_ = tree.Keys()   // []interface {}{1, 2, 3, 4, 5, 6} (in order)
//
//	tree.Remove(2) // 1->a, 3->c, 4->d, 5->e, 6->f (in order)
//	fmt.Println(tree)
//	//
//	//  RedBlackTree
//	//  │       ┌── 6
//	//  │   ┌── 5
//	//  └── 4
//	//      │   ┌── 3
//	//      └── 1
//
//	tree.Clear() // empty
//	tree.Empty() // true
//	tree.Size()  // 0
//
//	// Other:
//	tree.Left()     // gets the left-most (min) node
//	tree.Right()    // get the right-most (max) node
//	tree.Floor(1)   // get the floor node
//	tree.Ceiling(1) // get the ceiling node
//}
