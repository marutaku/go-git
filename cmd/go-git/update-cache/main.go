package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/marutaku/go-git/internal/cache"
	"github.com/marutaku/go-git/internal/env"
)

var activeCache cache.ActiveCache

func addCacheEntry(entry *cache.CacheEntry) error {
	existingEntryIndex := activeCache.FindCacheEntryIndex(entry)
	if existingEntryIndex != -1 {
		activeCache[existingEntryIndex] = entry
		return nil
	}
	activeCache = append(activeCache, entry)
	return nil
}

func addFileToCache(path string) error {
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
	entry, err := cache.NewCacheEntryFromFilePath(path, fileContent)
	if err != nil {
		return err
	}
	if activeCache.FindCacheEntryIndex(entry) != -1 {
		return addCacheEntry(entry)
	}
	err = entry.IndexFd(fileContent, stat)
	if err != nil {
		return err
	}
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
	srcIndexFilePath := fmt.Sprintf("%s/index.lock", env.GetSHA1FileDirectory())
	dstIndexFilePath := fmt.Sprintf("%s/index", env.GetSHA1FileDirectory())
	os.Rename(srcIndexFilePath, dstIndexFilePath)
}

func main() {
	var err error
	targetPaths := os.Args[1:]
	activeCache, err = cache.ReadCache()
	if err != nil {
		panic(err)
	}
	newIndexFile, err := os.OpenFile(fmt.Sprintf("%s/index.lock", env.GetSHA1FileDirectory()), os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		panic(err)
	}
	defer newIndexFile.Close()
	defer renameIndexFile()
	for _, path := range targetPaths {
		if !verifyPath(path) {
			continue
		}
		if err := addFileToCache(path); err != nil {
			panic(err)
		}
	}
	err = activeCache.WriteCache(newIndexFile)
	if err != nil {
		log.Fatal("unable to write cache")
	}

}
