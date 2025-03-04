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
