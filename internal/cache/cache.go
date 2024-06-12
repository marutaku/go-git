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
	CTime  CacheTime
	MTime  CacheTime
	STDev  uint32
	STIno  uint32
	STMode uint32
	STUid  uint32
	STGid  uint32
	STSize uint32
	Sha1   [20]byte
	Name   string
}

func (e *CacheEntry) Bytes() []byte {
	bytes := make([]byte, 0)
	bytes = append(bytes, byte(e.CTime.Sec))
	bytes = append(bytes, byte(e.CTime.NSec))
	bytes = append(bytes, byte(e.MTime.Sec))
	bytes = append(bytes, byte(e.MTime.NSec))
	bytes = append(bytes, byte(e.STDev))
	bytes = append(bytes, byte(e.STIno))
	bytes = append(bytes, byte(e.STMode))
	bytes = append(bytes, byte(e.STUid))
	bytes = append(bytes, byte(e.STGid))
	bytes = append(bytes, byte(e.STSize))
	bytes = append(bytes, e.Sha1[:]...)
	bytes = append(bytes, []byte(e.Name)...)
	return bytes
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

func NewCacheEntryFromFilePath(path string) (*CacheEntry, error) {
	fileStat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	ctime := &CacheTime{
		Sec:  uint32(fileStat.Sys().(*syscall.Stat_t).Ctimespec.Sec),
		NSec: uint32(fileStat.Sys().(*syscall.Stat_t).Ctimespec.Nsec),
	}
	mtime := &CacheTime{
		Sec:  uint32(fileStat.Sys().(*syscall.Stat_t).Mtimespec.Sec),
		NSec: uint32(fileStat.Sys().(*syscall.Stat_t).Mtimespec.Nsec),
	}
	entry := &CacheEntry{
		CTime:  *ctime,
		MTime:  *mtime,
		STDev:  uint32(fileStat.Sys().(*syscall.Stat_t).Dev),
		STIno:  uint32(fileStat.Sys().(*syscall.Stat_t).Ino),
		STMode: uint32(fileStat.Sys().(*syscall.Stat_t).Mode),
		STUid:  uint32(fileStat.Sys().(*syscall.Stat_t).Uid),
		STGid:  uint32(fileStat.Sys().(*syscall.Stat_t).Gid),
		STSize: uint32(fileStat.Size()),
		Name:   path,
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
