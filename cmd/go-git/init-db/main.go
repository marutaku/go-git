package main

import (
	"fmt"
	"os"

	"github.com/marutaku/go-git/internal"
)

func main() {
	if err := os.Mkdir(".dircache", 0700); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	sha1Dir := internal.DEFAULT_DB_ENVIRONMENT
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
