package main

import (
	"fmt"
	"log"
	"os"

	"github.com/marutaku/go-git/internal/env"
)

func main() {
	sha1Dir := env.GetSHA1FileDirectory()
	if err := os.Mkdir(sha1Dir, 0700); err != nil {
		log.Fatalf("error: %v\n", err)
	}
	if err := os.Mkdir(fmt.Sprintf("%s/objects", sha1Dir), 0700); err != nil {
		log.Fatalf("error: %v\n", err)
	}
	for i := 0; i < 256; i++ {
		dir := fmt.Sprintf("%s/objects/%02x", sha1Dir, i)
		if err := os.Mkdir(dir, 0700); err != nil {
			log.Fatalf("error: %v\n", err)
		}
	}
}
