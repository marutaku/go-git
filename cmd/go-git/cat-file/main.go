package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/marutaku/go-git/internal/buffer"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("cat-file: cat-file <sha1>")
	}
	sha1, err := hex.DecodeString(os.Args[1])
	if err != nil {
		log.Fatal("cat-file: cat-file <sha1>")
	}
	nodeType, buf, err := buffer.ReadSha1File(sha1)
	if err != nil {
		log.Fatal(err)
	}
	tmpfile, err := os.CreateTemp("", "temp_git_file_")
	if err != nil {
		log.Fatal(err)
	}
	defer tmpfile.Close()
	if _, err := tmpfile.Write(buf); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s: %s\n", tmpfile.Name(), nodeType)
}
