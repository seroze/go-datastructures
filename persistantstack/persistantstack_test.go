package persistantstack

import (
	"os"
	"testing"
)

func TestNaivePersistantStack(t *testing.T) {

	stack, err := NewNaivePersistantStack[int]("test_stack.db")

	if err != nil {
		t.Fatalf("Failed to create stack: %v", err)
	}

	defer stack.Close()

	defer os.Remove("test_stack.db") // Clean up after test

	stack.Push(10)
	stack.Push(20)
	stack.Push(242)
	stack.Push(224)

	if val, _ := stack.Pop(); val != 224 {
		t.Errorf("Expected 224, got %d", val)
	}

	if val, _ := stack.Pop(); val != 242 {
		t.Errorf("Expected 242, got %d", val)
	}
}
