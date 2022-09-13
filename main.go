package main

import (
	"os"
)

func main() {
	// Fisrt argument is for node's group
	arg_first := os.Args[1]
	// Second argument is for node's port
	arg_second := os.Args[2]

	if arg_first == "1" {
		NewNode(arg_second)
	}
}
