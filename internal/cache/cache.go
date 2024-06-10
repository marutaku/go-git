package cache

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"syscall"

	"github.com/marutaku/go-git/internal/buffer"
	"github.com/marutaku/go-git/internal/env"
)

type CacheTime struct {
	Sec  uint32
	NSec uint32
}

type CacheEntry struct {
	CTime   CacheTime
	MTime   CacheTime
	STDev   uint32
	STIno   uint32
	STMode  uint32
	STUid   uint32
	STGid   uint32
	STSize  uint32
	Sha1    [20]byte
	NameLen uint16
	Name    []byte
}

func IndexFd(nameLen int, entry *CacheEntry, fileContent string, stat fs.FileInfo) {
	contents := []byte(fmt.Sprintf("blob %d", stat.Size()))
	contents = append(contents, 0)
	contents = append(contents, []byte(fileContent)...)
	h := sha1.New()
	h.Write(contents)
	bs := h.Sum(nil)
	buffer.WriteSha1Buffer(bs, contents)
}

func updateCache(newIndexFile *os.File, path string) error {
	entry, err := NewCacheEntryFromFilePath(path)
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
	IndexFd(int(entry.NameLen), entry, string(fileContent), stat)
	return nil
}

func NewCacheEntryFromFilePath(path string) (*CacheEntry, error) {
	fileStat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	nameLen := len(path)
	ctime := &CacheTime{
		Sec:  uint32(fileStat.Sys().(*syscall.Stat_t).Ctimespec.Sec),
		NSec: uint32(fileStat.Sys().(*syscall.Stat_t).Ctimespec.Nsec),
	}
	mtime := &CacheTime{
		Sec:  uint32(fileStat.Sys().(*syscall.Stat_t).Mtimespec.Sec),
		NSec: uint32(fileStat.Sys().(*syscall.Stat_t).Mtimespec.Nsec),
	}
	entry := &CacheEntry{
		CTime:   *ctime,
		MTime:   *mtime,
		STDev:   uint32(fileStat.Sys().(*syscall.Stat_t).Dev),
		STIno:   uint32(fileStat.Sys().(*syscall.Stat_t).Ino),
		STMode:  uint32(fileStat.Sys().(*syscall.Stat_t).Mode),
		STUid:   uint32(fileStat.Sys().(*syscall.Stat_t).Uid),
		STGid:   uint32(fileStat.Sys().(*syscall.Stat_t).Gid),
		STSize:  uint32(fileStat.Size()),
		NameLen: uint16(nameLen),
	}
	return entry, nil
}

func ReadCache() (int, error) {
	sha1FileDir := env.GetSHA1FileDirectory()
	if _, err := os.Stat(sha1FileDir); os.IsExist(err) {
		return 0, errors.New("SHA1 file directory not found")
	}
	// TODO: Implement the rest of the function
	return 0, nil
}
