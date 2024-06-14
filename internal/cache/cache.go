package cache

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"syscall"

	objectBuffer "github.com/marutaku/go-git/internal/buffer"
	"github.com/marutaku/go-git/internal/cache/cachetime"
	"github.com/marutaku/go-git/internal/env"
)

type CacheEntry struct {
	CTime   cachetime.CacheTime
	MTime   cachetime.CacheTime
	STDev   uint32
	STIno   uint32
	STMode  uint32
	STUid   uint32
	STGid   uint32
	STSize  uint32
	Sha1    []byte
	NameLen uint16
	Name    string
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
	bytes = append(bytes, byte(e.NameLen))
	bytes = append(bytes, []byte(e.Name)...)
	return bytes
}

func (e *CacheEntry) IndexFd(fileContent string, stat fs.FileInfo) error {
	contents := []byte(fmt.Sprintf("blob %d", stat.Size()))
	contents = append(contents, 0)
	contents = append(contents, []byte(fileContent)...)
	var buffer bytes.Buffer
	zWriter := zlib.NewWriter(&buffer)
	zWriter, err := zlib.NewWriterLevel(zWriter, zlib.BestCompression)
	if err != nil {
		return err
	}
	zWriter.Write(contents)
	h := sha1.New()
	h.Write(contents)
	sha1Bytes := h.Sum(nil)
	e.Sha1 = sha1Bytes
	objectBuffer.WriteSha1Buffer(sha1Bytes, buffer.Bytes())
	return nil
}

func IndexFd(nameLen int, entry *CacheEntry, fileContent string, stat fs.FileInfo) error {
	contents := []byte(fmt.Sprintf("blob %d", stat.Size()))
	contents = append(contents, 0)
	contents = append(contents, []byte(fileContent)...)
	var buffer bytes.Buffer
	zWriter := zlib.NewWriter(&buffer)
	zWriter, err := zlib.NewWriterLevel(zWriter, zlib.BestCompression)
	if err != nil {
		return err
	}
	zWriter.Write(contents)
	h := sha1.New()
	h.Write(contents)
	sha1Bytes := h.Sum(nil)
	objectBuffer.WriteSha1Buffer(sha1Bytes, buffer.Bytes())
	return nil
}

func NewCacheEntryFromFilePath(path string) (*CacheEntry, error) {
	fileStat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	// https://github.com/golang/go/issues/29393
	ctime := cachetime.NewCTimeFromStat(fileStat)
	mtime := cachetime.NewMTimeFromStat(fileStat)
	entry := &CacheEntry{
		CTime:   *ctime,
		MTime:   *mtime,
		STDev:   uint32(fileStat.Sys().(*syscall.Stat_t).Dev),
		STIno:   uint32(fileStat.Sys().(*syscall.Stat_t).Ino),
		STMode:  uint32(fileStat.Sys().(*syscall.Stat_t).Mode),
		STUid:   uint32(fileStat.Sys().(*syscall.Stat_t).Uid),
		STGid:   uint32(fileStat.Sys().(*syscall.Stat_t).Gid),
		STSize:  uint32(fileStat.Size()),
		NameLen: uint16(len(path)),
		Name:    path,
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
