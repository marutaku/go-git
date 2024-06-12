package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/marutaku/go-git/internal/cache"
)

type ActiveCache []*cache.CacheEntry

func (ac ActiveCache) findCacheEntryIndex(path string) int {
	for index, entry := range ac {
		if entry.Name == path {
			return index
		}
	}
	return -1
}

func (ac ActiveCache) writeCache(file *os.File) error {
	// SHA1ハッシュとる箇所の自信がない
	header := cache.NewCacheHeader(1, ac)
	headerBytes := header.Bytes()
	file.Write(headerBytes)
	for _, entry := range ac {
		entryBytes := entry.Bytes()
		file.Write(entryBytes)
	}
	return nil
}

var activeCache ActiveCache

func addCacheEntry(entry *cache.CacheEntry) error {
	existingEntryIndex := activeCache.findCacheEntryIndex(entry.Name)
	if existingEntryIndex != -1 {
		activeCache[existingEntryIndex] = entry
		return nil
	}
	activeCache = append(activeCache, entry)
	return nil
}

func addFileToCache(path string) error {
	entry, err := cache.NewCacheEntryFromFilePath(path)
	if err != nil {
		return err
	}
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	fileContent := make([]byte, stat.Size())
	_, err = file.Read(fileContent)
	if err != nil {
		return err
	}
	cache.IndexFd(len(entry.Name), entry, string(fileContent), stat)
	return addCacheEntry(entry)
}

func verifyPath(path string) bool {
	if strings.Contains(path, "..") || strings.Contains(path, "//") || strings.HasSuffix(path, "/") {
		return false
	}
	if filepath.Base(path)[0] == '.' {
		return false
	}
	return true
}

func renameIndexFile() {
	os.Rename(".dircache/index.lock", ".dircache/index")
}

func main() {
	targetPaths := os.Args[1:]
	entries, err := cache.ReadCache()
	if err != nil {
		panic(err)
	}
	if entries < 0 {
		log.Fatal("cache corrupted")
	}

	newIndexFile, err := os.Create(".dircache/index.lock")
	if err != nil {
		log.Fatal("unable to create new cache file")
	}
	defer newIndexFile.Close()
	defer renameIndexFile()
	for _, path := range targetPaths {
		if !verifyPath(path) {
			fmt.Printf("Ignoring path %s\n", path)
			continue
		}
		if err := addFileToCache(path); err != nil {
			log.Fatalf("Unable to add %s to database\n", path)
		}
	}
	err = activeCache.writeCache(newIndexFile)
	if err != nil {
		log.Fatal("unable to write cache")
	}

}
