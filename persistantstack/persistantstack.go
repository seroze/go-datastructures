package persistantstack

import (
	"bufio"
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"sync"
)

/***
It's a basic version that just appends the elements to the end of a file
and when we need to pop it goes through the entire file to figure out where the last element is

***/

type NaivePersistantStack[T any] struct {
	file   *os.File
	writer *bufio.Writer
	mutex  sync.Mutex
}

func NewNaivePersistantStack[T any](filename string) (*NaivePersistantStack[T], error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)

	if err != nil {
		return nil, err
	}

	return &NaivePersistantStack[T]{
		file:   file,
		writer: bufio.NewWriter(file),
	}, nil
}

func (nps *NaivePersistantStack[T]) Push(value T) error {

	nps.mutex.Lock()
	defer nps.mutex.Unlock()

	enc := gob.NewEncoder(nps.writer)

	if err := enc.Encode(value); err != nil {
		return err
	}

	return nps.writer.Flush()

}

func (nps *NaivePersistantStack[T]) Pop() (T, error) {
	nps.mutex.Lock()
	defer nps.mutex.Unlock()

	var zero T

	info, err := nps.file.Stat()
	if err != nil {
		return zero, err
	}

	if info.Size() == 0 {
		return zero, errors.New("stack underflow: no elements left")
	}

	nps.file.Seek(0, 0) // Move to the beginning
	dec := gob.NewDecoder(nps.file)

	var lastElement T
	var lastPos int64 = -1

	for {
		curPos, _ := nps.file.Seek(0, os.SEEK_CUR) //get current offset before reading
		fmt.Println(curPos, " current Pos")
		var value T
		err := dec.Decode(&value)
		if err != nil {
			break //Stop reading file when we hit EOF
		}
		lastElement = value
		lastPos = curPos
	}

	// If we read at least one value, truncate at prevPos
	if lastPos > 0 {
		fmt.Printf("🛑 Truncating at position: %d\n", lastPos)
		err = nps.file.Truncate(lastPos)

		if err != nil {
			return zero, err
		}

		// **Ensure the cursor is at the end after truncation**
		nps.file.Seek(0, os.SEEK_END)
		nps.file.Sync() //force flush changes

		info, _ = nps.file.Stat()
		fmt.Println("✅ File Size After Pop:", info.Size())
	}

	return lastElement, nil
}

func (nps *NaivePersistantStack[T]) Close() error {
	nps.mutex.Lock()
	defer nps.mutex.Unlock()

	return nps.file.Close()
}
