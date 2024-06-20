package main

import (
	"bytes"
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
	defer file.Close()
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
	if index := activeCache.FindCacheEntryIndex(entry); index != -1 {
		if bytes.Equal(entry.Sha1, activeCache[index].Sha1) {
			// 全く同じであれば何もしない
			return nil
		}
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
	tmpIndexFilePath := fmt.Sprintf("%s/index.lock", env.GetSHA1FileDirectory())
	newIndexFile, err := os.OpenFile(tmpIndexFilePath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		panic(err)
	}
	defer newIndexFile.Close()
	for _, path := range targetPaths {
		if !verifyPath(path) {
			continue
		}
		if err := addFileToCache(path); err != nil {
			os.Remove(tmpIndexFilePath)
			log.Fatal("unable to add file to cache: ", err)
		}
	}
	err = activeCache.WriteCache(newIndexFile)
	if err != nil {
		os.Remove(tmpIndexFilePath)
		log.Fatal("unable to write cache: ", err)
	}
	renameIndexFile()

}
