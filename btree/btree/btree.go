package btree

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

func SayHello(name string) {
	fmt.Println("Hello", name)
}

func init() {
	SayHello("Module initialized")
}

type BTreeNode[T constraints.Ordered] struct {
	keys     []T
	children []*BTreeNode[T]
	t        int
	leaf     bool
}

func NewBTreeNode[T constraints.Ordered](t int, isLeaf bool) *BTreeNode[T] {
	return &BTreeNode[T]{
		keys:     make([]T, 0, 2*t-1),
		children: make([]*BTreeNode[T], 0, 2*t),
		t:        t,
		leaf:     isLeaf,
	}
}

type BTree[T constraints.Ordered] struct {
	root *BTreeNode[T]
	t    int
}

// NewBtree
func NewBTree[T constraints.Ordered](t int) *BTree[T] {
	return &BTree[T]{root: NewBTreeNode[T](t, true), t: t}
}

func (tree *BTree[T]) Search(key T) bool {
	root := tree.root
	node, _ := root.search(key)
	if node == nil {
		return false
	} else {
		return true
	}
}

func (node *BTreeNode[T]) search(target T) (*BTreeNode[T], int) {
	i := 0
	for i < len(node.keys) && target > node.keys[i] {
		i++
	}
	if i < len(node.keys) && target == node.keys[i] {
		return node, i
	} else if node.leaf == true {
		return nil, -1
	} else if i < len(node.children) {
		return node.children[i].search(target)
	} else {
		return nil, -1
	}
}

func insertAt[T any](slice []T, index int, element T) []T {
	updatedSlice := make([]T, 0, len(slice)+1)            // Preallocate capacity
	updatedSlice = append(updatedSlice, slice[:index]...) // First part
	updatedSlice = append(updatedSlice, element)          // Insert element
	updatedSlice = append(updatedSlice, slice[index:]...) // Second part
	return updatedSlice
}

func (node *BTreeNode[T]) splitChild(i int) {
	// the child which we want to split
	y := node.children[i]
	t := node.t

	// assert that this child is indeed full
	// assert that current node is not full cause we insert the median into it

	z := NewBTreeNode[T](node.t, y.leaf)
	// insert z into children at i
	node.children = insertAt(node.children, i+1, z)
	// insert median of y into keys at i
	node.keys = insertAt(node.keys, i, y.keys[t-1])

	mid := t - 1

	//move right half to z
	z.keys = y.keys[mid+1:]
	//keep the left half in y except the median element
	y.keys = y.keys[:mid]

	// if y is not a leaf then distribute y's children accordingly
	if y.leaf == false {
		// think through these cases more clearly
		z.children = y.children[mid+1:]
		y.children = y.children[0 : mid+1]
	}
}

func (tree *BTree[T]) Insert(key T) {
	root := tree.root
	if len(root.keys) == 2*tree.t-1 { // If root is full, create a new root
		newNode := NewBTreeNode[T](tree.t, false)
		newNode.children = append(newNode.children, root)
		newNode.splitChild(0)
		newNode.insertNonFull(key)
		tree.root = newNode //  Update root
	} else {
		root.insertNonFull(key)
	}
}

func (node *BTreeNode[T]) insertNonFull(key T) {

	// fmt.Println("Inserting in non-full for ", key, node.keys)
	i := len(node.keys) - 1
	t := node.t
	if node.leaf == true {
		var zeroValue T
		node.keys = append(node.keys, zeroValue)
		// fmt.Println(node.keys, " in insertNonFull")
		for i >= 0 && key < node.keys[i] {
			//move ith key to i+1th place
			node.keys[i+1] = node.keys[i]
			i--
		}
		node.keys[i+1] = key
	} else {
		//if it's not leaf
		// check for the right child
		for i >= 0 && key < node.keys[i] {
			i--
		}
		i++
		// if the child node is almost full, split it
		if len(node.children[i].keys) == (2*t)-1 {
			node.splitChild(i)
			if key > node.keys[i] {
				i++
			}
		}
		node.children[i].insertNonFull(key)
	}
}

func (tree *BTree[T]) Delete(target T) {
	root := tree.root
	t := root.t
	// first search for the key
	node, index := root.search(target)
	if node == nil {
		// do nothing
		return
	}
	// Case 1: if it's a leaf
	if node.leaf {
		//just delete it, but make sure it has >=t leaves
		if len(node.keys) >= t {
			j := index
			for j+1 < len(node.keys) {
				//swap j with it's immediate right
				node.keys[j] = node.keys[j+1]
				j++
			}
			// pop the last element
			node.keys = node.keys[:len(node.keys)-1]
		}
	}

}

// PrintBTree prints the structure of the B-Tree
func (tree *BTree[T]) PrintBTree() {
	if tree.root == nil {
		fmt.Println("Tree is empty")
		return
	}
	tree.root.printNode(0) // Start with root at depth 0
}

// printNode recursively prints the node and its children
func (node *BTreeNode[T]) printNode(depth int) {
	if node == nil {
		return
	}

	// Indentation based on depth
	for i := 0; i < depth; i++ {
		fmt.Print("  ")
	}
	fmt.Println(node.keys)

	// Print children recursively
	for _, child := range node.children {
		child.printNode(depth + 1)
	}
}
