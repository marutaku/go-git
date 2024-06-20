package main

import (
	"fmt"
	"log"
	"os"

	"github.com/marutaku/go-git/internal/buffer"
	"github.com/marutaku/go-git/internal/cache"
)

var ORIG_OFFSET = 40

func checkValidSha1(sha1Hash []byte) bool {
	fileName := buffer.GetSha1FileName(sha1Hash)
	_, err := os.Stat(fileName)
	return err == nil
}

func main() {
	entries, err := cache.ReadCache()
	if err != nil {
		panic(err)
	}
	if len(entries) == 0 {
		log.Fatal("No file-cache to create a tree of \n")
	}
	size := len(entries)*40 + 400
	offset := ORIG_OFFSET
	treeBuffer := make([]byte, size)
	for _, entry := range entries {
		if !checkValidSha1(entry.Sha1) {
			log.Fatalf("Invalid sha1: %x\n", entry.Sha1)
		}
		requiredSpace := offset + int(entry.NameLen) + 60
		if requiredSpace > size {
			size = ((requiredSpace) + 16) * 3 / 2
			treeBuffer = append(treeBuffer, make([]byte, size-len(treeBuffer))...)
		}
		copy(treeBuffer[offset:], []byte(fmt.Sprintf("%o %s", entry.STMode, entry.Name)))
		offset += requiredSpace
		treeBuffer[offset+1] = 0
		offset++
		copy(treeBuffer[offset:], entry.Sha1)
		offset += 20
	}
	i := buffer.PrependInteger(treeBuffer, offset-ORIG_OFFSET, ORIG_OFFSET)
	i -= 5
	copy(treeBuffer[i:], []byte("tree "))
	offset -= i
	err = cache.WriteSha1File(treeBuffer[i:])
	if err != nil {
		log.Fatal("Failed to write tree: ", err)
	}
}
