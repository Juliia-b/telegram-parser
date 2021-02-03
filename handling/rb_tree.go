package handling

import (
	"fmt"
	rbt "github.com/emirpasic/gods/trees/redblacktree"
	"sync"
)

//    --------------------------------------------------------------------------------
//                                     METHODS
//    --------------------------------------------------------------------------------

// CreateRBTree creates a new empty tree and lists it in the tree list
func (h *HashRing) CreateRBTree(nodeAddr string) (tree *rbt.Tree) {
	tree = rbt.NewWithStringComparator()
	addTreeToMap(h, nodeAddr, tree)
	return tree
}

// DeleteRbTree removes a tree from the tree list
func (h *HashRing) DeleteRbTree(nodeAddr string) {
	delete(h.Trees, nodeAddr)
}

// getRbTree returns a red-black tree for a node with a given address
func (h *HashRing) getRbTree(nodeAddr string) (tree *rbt.Tree, ok bool) {
	tree, ok = getTreeFromMap(h, nodeAddr)
	return tree, ok
}

// AddValToTree adds a new value to the red-black tree.
// @param nodeAddr is the tree identifier.
// @param value is a unique key of the telegram message.
func (h *HashRing) AddValToTree(nodeAddr string, value string) {
	tree, ok := h.getRbTree(nodeAddr)
	if !ok {
		// there is no tree with that name. Need to create
		tree = h.CreateRBTree(nodeAddr)
	}

	tree.Put(value, value)
}

// GetTreeValues gets all values in the tree
func (h *HashRing) GetTreeValues(nodeAddr string) (values []string) {
	var res []string

	tree, ok := h.getRbTree(nodeAddr)
	if !ok {
		return res
	}

	for _, v := range tree.Values() {
		res = append(res, fmt.Sprintln(v))
	}

	return res
}

//    --------------------------------------------------------------------------------
//                                     HELPERS
//    --------------------------------------------------------------------------------

func addTreeToMap(h *HashRing, nodeAddr string, tree *rbt.Tree) {
	mu := sync.RWMutex{}
	mu.Lock()
	h.Trees[nodeAddr] = tree
	mu.Unlock()
}

func getTreeFromMap(h *HashRing, nodeAddr string) (tree *rbt.Tree, exist bool) {
	mu := sync.RWMutex{}
	mu.RLock()
	tree, exist = h.Trees[nodeAddr]
	mu.RUnlock()
	return tree, exist
}

//    --------------------------------------------------------------------------------
//                                        EXTRA
//    --------------------------------------------------------------------------------

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
