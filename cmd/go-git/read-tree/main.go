package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/marutaku/go-git/internal/buffer"
)

func unpack(sha1 []byte) error {
	nodeType, byteBuffer, err := buffer.ReadSha1File(sha1)
	if err != nil {
		return err
	}
	if nodeType != "tree" {
		return fmt.Errorf("invalid node type: %s", nodeType)
	}
	fileInfos := bytes.Split(byteBuffer, []byte(" "))
	for _, fileInfo := range fileInfos {
		var mode int
		var name string
		fmt.Println(string(fileInfo))
		splitted := bytes.Split(fileInfo, []byte("\x00"))
		fmt.Sscanf(string(splitted[0]), "%o %s (%x)\n", &mode, &name, splitted[1])
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
