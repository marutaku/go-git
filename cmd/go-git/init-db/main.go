package main

import (
	"fmt"
	"os"

	"github.com/marutaku/go-git/internal/env"
)

func main() {
	sha1Dir := env.GetSHA1FileDirectory()
	if err := os.Mkdir(sha1Dir, 0700); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	for i := 0; i < 256; i++ {
		dir := fmt.Sprintf("%s/%02x", sha1Dir, i)
		if err := os.Mkdir(dir, 0700); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}
}
