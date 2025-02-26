package persistantstack 

import "os"
import "testing" 

func TestNaivePersistantStack(t *testing.T) {

	stack, err := NewNaivePersistantStack[int]("test_stack.db")

	if err!=nil{
		t.Fatalf("Failed to create stack: %v", err)
	}

	defer stack.Close() 

	defer os.Remove("test_stack.db") // Clean up after test 

	stack.Push(10)
	stack.Push(242)
	stack.Push(224)

	if val, _ := stack.Pop(); val!=30{
		t.Errorf("Expected 30, got %d", val)
	}

	if val,_ := stack.Pop(); val!=20{
		t.Errorf("Expected 20, got %d", val)
	}
}
