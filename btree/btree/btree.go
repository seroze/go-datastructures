package btree

import "fmt"

func SayHello(name string) {
	fmt.Println("Hello", name)
}

type BTreeNode struct {
	keys     []int
	children []*BTreeNode
	leaf     bool
}

func NewBTreeNode() *BTreeNode {
	return &BTreeNode{
		keys: make([]int, 0, 2*t-1),
		children: make([]*BTreeNode, 0, 2*t),
		leaf: true,
	}
}

type BTree sturct {
	root *BTreeNode
	t int
}

//NewBtree
func NewBTree(t int){
	return & Btree{root: NewBTreeNode(),  t: t}
}

func (node* BTreeNode) splitChild()
