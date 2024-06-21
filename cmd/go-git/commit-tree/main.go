package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("commit-tree <sha1> [-p <sha1>]* < changelog")
	}
	
}
