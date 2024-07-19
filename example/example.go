package main

import (
	"fmt"
)

func main() {
	err := doSomething()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Success!")
}

func doSomething() error {
	// Simulate some work
	return nil
}
