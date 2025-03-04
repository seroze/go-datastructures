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
	// what's the point of splitting if the root is full ?
	// cause it's to create non-full node at the top and at the
	// same time it creates non-full child's
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
		// we recursively keep splitting the child if it's full
		// and at the same time the parent is definetly non-full
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
	if tree.root == nil {
		return
	}

	tree.root.delete(target)

	// if the root node has no keys left
	// you'll understand this better once you go through the
	// rest of the explanation
	if len(tree.root.keys) == 0 {
		if !tree.root.leaf {
			tree.root = tree.root.children[0]
		} else {
			tree.root = nil
		}
	}
}

func (node *BTreeNode[T]) delete(key T) {
	t := node.t
	i := 0
	for i < len(node.keys) && key > node.keys[i] {
		i++
	}

	//Case 1: if the key is present inside this node
	if i < len(node.keys) && node.keys[i] == key {
		// Case 1a: it's a leaf node
		if node.leaf {
			node.deleteFromLeaf(i)
		} else {
			// Case 1b: it's an internal node
			node.deleteFromInternalNode(i)
		}
	} else {
		// Case 2: key is not present in this node
		if node.leaf {
			// key is not present in this leaf node
			return
		}

		// Case 2a: check if child doesn't have enough keys
		if len(node.children[i].keys) < t {
			//children[i] is guaranteed to exist at this point
			// borrow from left or right but fill it
			node.fillChild(i)
		}

		// If the last child was merged from above then recurse on the (i-1)th child
		if i > len(node.keys) {
			node.children[i-1].delete(key)
		} else {
			node.children[i].delete(key)
		}
	}
}

func (node *BTreeNode[T]) deleteFromLeaf(index int) {
	// but what if this node has only t-1 keys ?
	// deleteFromInternal takes care of this
	node.keys = append(node.keys[:index], node.keys[index+1:]...)
}

// delets node.keys[index] from current node, it's assumed we have atleast t keys here
func (node *BTreeNode[T]) deleteFromInternalNode(index int) {
	key := node.keys[index]
	t := node.t

	// Case 1: if left child has atleast t keys, replace with predecessor
	if len(node.children[index].keys) >= t {
		pred := node.getPredecessor(index)
		node.keys[index] = pred
		node.children[index].delete(pred)
	} else if len(node.children[index+1].keys) >= t {
		// Case 2: if right child has atleast t keys, replace with successor
		succ := node.getSuccessor(index)
		node.keys[index] = succ
		node.children[index+1].delete(succ)
	} else {
		// Case 3: merge the key and the right child into the left child
		// for merge to work current node should have atleast t keys
		node.mergeChildren(index)
		node.children[index].delete(key)
	}
}

// getPredecessor finds the predecessor of a key in the subtree
func (node *BTreeNode[T]) getPredecessor(index int) T {
	curr := node.children[index]
	// keep moving along the right most child
	for !curr.leaf {
		curr = curr.children[len(curr.children)-1]
	}
	return curr.keys[len(curr.keys)-1]
}

func (node *BTreeNode[T]) getSuccessor(index int) T {
	curr := node.children[index+1]

	// keep moving along the leftmost child
	for !curr.leaf {
		curr = curr.children[0]
	}
	return curr.keys[0]
}

func (node *BTreeNode[T]) fillChild(index int) {
	t := node.t
	if index != 0 && len(node.children[index-1].keys) >= t {
		// Borrow from the left sibling
		node.borrowFromLeft(index)
	} else if index+1 < len(node.children) && len(node.children[index+1].keys) >= t {
		// Borrow from the right sibling
		node.borrowFromRight(index)
	} else {
		// merge with a sibling
		// i have no clue what this is
		if index != len(node.keys) {
			node.mergeChildren(index)
		} else {
			node.mergeChildren(index - 1)
		}
	}
}

func (node *BTreeNode[T]) borrowFromLeft(index int) {
	child := node.children[index]
	leftSibling := node.children[index-1]

	// move key from parent to the child
	child.keys = append([]T{node.keys[index-1]}, child.keys...)
	// replace the parent's key with left sibling's rightmost key
	node.keys[index-1] = leftSibling.keys[len(leftSibling.keys)-1]

	// move the right most child of left sibling as left child to interested node
	if !leftSibling.leaf {
		child.children = append([]*BTreeNode[T]{leftSibling.children[len(leftSibling.keys)-1]}, child.children...)
		leftSibling.children = leftSibling.children[:len(leftSibling.children)-1]
	}

	// remove the borrowed key from the left sibling
	leftSibling.keys = leftSibling.keys[:len(leftSibling.keys)-1]
}

// borrowFromRight borrows a key from the right sibling
func (node *BTreeNode[T]) borrowFromRight(index int) {
	child := node.children[index]
	rightSibling := node.children[index+1]

	// Move a key from the parent to the child
	child.keys = append(child.keys, node.keys[index])
	node.keys[index] = rightSibling.keys[0]

	// Move the first child of the right sibling to the child
	if !rightSibling.leaf {
		child.children = append(child.children, rightSibling.children[0])
		rightSibling.children = rightSibling.children[1:]
	}

	// Remove the borrowed key from the right sibling
	rightSibling.keys = rightSibling.keys[1:]
}

// one last function
// merges child with it's right sibling and delete's the right sibling after merging with parent
func (node *BTreeNode[T]) mergeChildren(index int) {
	child := node.children[index]
	rightSibling := node.children[index+1]

	// move the key from the parent to the children
	child.keys = append(child.keys, node.keys[index])
	child.keys = append(child.keys, rightSibling.keys...)

	// move the children of the rightSibling to the child
	if !child.leaf {
		child.children = append(child.children, rightSibling.children...)
	}

	// remove the key and rightSibling from the parent
	node.keys = append(node.keys[:index], node.keys[index+1:]...)
	// does it overflow
	node.children = append(node.children[:index+1], node.children[index+2:]...)
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
