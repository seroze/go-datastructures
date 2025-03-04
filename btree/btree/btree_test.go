package btree

import (
	"fmt"
	"testing"
)

// TestBTreeInsert tests insertion of keys into the B-tree.
func TestBTreeInsert(t *testing.T) {
	// tree := NewBTreeting t=2 for easier testing
	tree := NewBTree[int](2)

	keysToInsert := []int{10, 20, 5, 6, 12, 30, 7, 17}

	for _, key := range keysToInsert {

		tree.Insert(key)
		fmt.Println(" after inserting ", key)
		tree.PrintBTree()
	}

	// Ensure all inserted keys are found
	for _, key := range keysToInsert {
		fmt.Println("searching for", key)
		if !tree.Search(key) {
			t.Errorf("Expected key %d to be found in the BTree, but it was not.", key)
		}
	}

	// Ensure searching for a non-existent key returns false
	nonExistentKeys := []int{100, -5, 42}
	for _, key := range nonExistentKeys {
		if tree.Search(key) {
			t.Errorf("Expected key %d to NOT be found in the BTree, but it was.", key)
		}
	}
}

// TestBTreeInsertDuplicates ensures duplicates are handled properly.
func TestBTreeInsertDuplicates(t *testing.T) {
	tree := NewBTree[int](2)

	// uplicate keys
	tree.Insert(10)
	tree.Insert(10)

	// Check if 10 is found
	if !tree.Search(10) {
		t.Errorf("Expected key 10 to be found in the BTree.")
	}

	// Since we're not explicitly handling duplicates, it should just be found once
	// (Structural validation would require more in-depth tests)
}

func TestBTreeDeletion(t *testing.T) {
	fmt.Println("Running B-Tree Deletion Tests...")

	tDegree := 2 // Minimum degree
	btree := BTree[int]{root: NewBTreeNode[int](tDegree, true), t: tDegree}

	// Insert elements
	elements := []int{10, 20, 5, 6, 12, 30, 7, 17}
	for _, el := range elements {
		btree.Insert(el)
	}

	fmt.Println("Initial B-Tree:")
	btree.PrintBTree()

	// Case 1: Deleting a leaf node
	t.Run("DeleteLeafNode", func(t *testing.T) {
		fmt.Println("\nDeleting leaf node 6...")
		btree.Delete(6)
		btree.PrintBTree()
		if btree.Search(6) {
			t.Errorf("Key 6 should have been deleted but is still present.")
		}
	})

	// Case 2a: Deleting an internal node where left child has enough keys
	t.Run("DeleteInternalNodeLeft", func(t *testing.T) {
		fmt.Println("\nDeleting internal node 10...")
		btree.Delete(10)
		btree.PrintBTree()
		if btree.Search(10) {
			t.Errorf("Key 10 should have been deleted but is still present.")
		}
	})

	// Case 2b: Deleting an internal node where right child has enough keys
	t.Run("DeleteInternalNodeRight", func(t *testing.T) {
		fmt.Println("\nDeleting internal node 12...")
		btree.Delete(12)
		btree.PrintBTree()
		if btree.Search(12) {
			t.Errorf("Key 12 should have been deleted but is still present.")
		}
	})

	// Case 2c: Deleting an internal node where both children have t-1 keys
	t.Run("DeleteInternalNodeBothTMinus1", func(t *testing.T) {
		fmt.Println("\nDeleting internal node 20...")
		btree.Delete(20)
		btree.PrintBTree()
		if btree.Search(20) {
			t.Errorf("Key 20 should have been deleted but is still present.")
		}
	})

	// Case 3: Deleting from a node that requires merging
	t.Run("DeleteNodeCausesMerge", func(t *testing.T) {
		fmt.Println("\nDeleting node 30 (causes merge)...")
		btree.Delete(30)
		btree.PrintBTree()
		if btree.Search(30) {
			t.Errorf("Key 30 should have been deleted but is still present.")
		}
	})
}
