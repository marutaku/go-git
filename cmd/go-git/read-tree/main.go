package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/marutaku/go-git/internal/objects"
)

func unpack(sha1 []byte) error {
	nodeType, byteBuffer, err := objects.ReadSha1File(sha1)
	if err != nil {
		return err
	}
	if nodeType != "tree" {
		return fmt.Errorf("invalid node type: %s", nodeType)
	}
	size := len(byteBuffer)
	offset := 0
	for offset < size {
		var mode int
		var name string
		var sha1 []byte
		nullByteIndex := bytes.IndexByte(byteBuffer[offset:], 0)
		fileInfo := string(byteBuffer[offset : offset+nullByteIndex])
		fmt.Sscanf(fileInfo, "%o %s", &mode, &name)
		offset += nullByteIndex + 1
		sha1 = byteBuffer[offset : offset+20]
		fmt.Printf("%o %s (%x)\n", mode, name, sha1)
		offset += 20
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: read-tree <key>")
		os.Exit(1)
	}
	sha1, err := hex.DecodeString(os.Args[1])
	if err != nil {
		fmt.Println("Invalid sha1")
		os.Exit(1)
	}
	err = unpack(sha1)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
