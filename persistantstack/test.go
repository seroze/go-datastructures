package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	file, err := os.Open("data.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	var p Person
	for {
		pos, _ := file.Seek(0, os.SEEK_CUR) // Get current position before reading
		fmt.Println("Current file offset:", pos)

		err := decoder.Decode(&p)
		if err != nil {
			break
		}

		// Workaround: Re-seek to actual file position
		newPos, _ := file.Seek(0, os.SEEK_CUR)
		fmt.Println("New file offset:", newPos)
		fmt.Printf("Decoded: %+v\n", p)
	}
}
