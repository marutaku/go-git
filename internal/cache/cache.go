package cache

import (
	"crypto/sha1"
	"encoding/binary"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"syscall"

	objectBuffer "github.com/marutaku/go-git/internal/buffer"
	"github.com/marutaku/go-git/internal/cache/cachetime"
	"github.com/marutaku/go-git/internal/env"
	"github.com/marutaku/go-git/internal/utils"
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
	bytes = binary.LittleEndian.AppendUint32(bytes, e.CTime.Sec)
	bytes = binary.LittleEndian.AppendUint32(bytes, e.CTime.NSec)
	bytes = binary.LittleEndian.AppendUint32(bytes, e.MTime.Sec)
	bytes = binary.LittleEndian.AppendUint32(bytes, e.MTime.NSec)
	bytes = binary.LittleEndian.AppendUint32(bytes, e.STDev)
	bytes = binary.LittleEndian.AppendUint32(bytes, e.STIno)
	bytes = binary.LittleEndian.AppendUint32(bytes, e.STMode)
	bytes = binary.LittleEndian.AppendUint32(bytes, e.STUid)
	bytes = binary.LittleEndian.AppendUint32(bytes, e.STGid)
	bytes = binary.LittleEndian.AppendUint32(bytes, e.STSize)
	bytes = append(bytes, e.Sha1[:]...)
	bytes = binary.LittleEndian.AppendUint16(bytes, e.NameLen)
	bytes = append(bytes, []byte(e.Name)...)
	return bytes
}

func calculateSha1Hash(stat fs.FileInfo, fileContent []byte) ([]byte, error) {
	contents := []byte(fmt.Sprintf("blob %d", uint32(stat.Size())))
	contents = append(contents, 0)
	contents = append(contents, []byte(fileContent)...)
	compressed, err := utils.Compress(contents)
	if err != nil {
		return nil, err
	}
	h := sha1.New()
	h.Write(compressed)
	sha1Bytes := h.Sum(nil)
	return sha1Bytes, nil
}

func (e *CacheEntry) IndexFd(fileContent []byte, stat fs.FileInfo) error {
	contents := []byte(fmt.Sprintf("blob %d", uint32(stat.Size())))
	contents = append(contents, 0)
	contents = append(contents, fileContent...)
	compressed, err := utils.Compress(contents)
	if err != nil {
		return err
	}
	return objectBuffer.WriteSha1Buffer(e.Sha1, compressed)
}

func NewCacheEntryFromFilePath(path string, fileContents []byte) (*CacheEntry, error) {
	fileStat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	sha1, err := calculateSha1Hash(fileStat, fileContents)
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
		Sha1:    sha1,
	}
	return entry, nil
}

func NewCacheEntryFromBytes(indexFileBytes []byte) (*CacheEntry, uint32) {
	entry := &CacheEntry{}
	entry.CTime.Sec = binary.LittleEndian.Uint32(indexFileBytes[:4])
	entry.CTime.NSec = binary.LittleEndian.Uint32(indexFileBytes[4:8])
	entry.MTime.Sec = binary.LittleEndian.Uint32(indexFileBytes[8:12])
	entry.MTime.NSec = binary.LittleEndian.Uint32(indexFileBytes[12:16])
	entry.STDev = binary.LittleEndian.Uint32(indexFileBytes[16:20])
	entry.STIno = binary.LittleEndian.Uint32(indexFileBytes[20:24])
	entry.STMode = binary.LittleEndian.Uint32(indexFileBytes[24:28])
	entry.STUid = binary.LittleEndian.Uint32(indexFileBytes[28:32])
	entry.STGid = binary.LittleEndian.Uint32(indexFileBytes[32:36])
	entry.STSize = binary.LittleEndian.Uint32(indexFileBytes[36:40])
	entry.Sha1 = indexFileBytes[40:60]
	entry.NameLen = binary.LittleEndian.Uint16(indexFileBytes[60:62])
	entry.Name = string(indexFileBytes[62 : 62+entry.NameLen])
	return entry, uint32(62 + entry.NameLen)
}

func ReadCache() (ActiveCache, error) {
	sha1FileDir := env.GetSHA1FileDirectory()
	if _, err := os.Stat(sha1FileDir); os.IsExist(err) {
		return nil, errors.New("SHA1 file directory not found")
	}
	if _, err := os.Stat(fmt.Sprintf("%s/index", sha1FileDir)); os.IsNotExist(err) {
		return ActiveCache{}, nil
	}
	bytes, err := os.ReadFile(fmt.Sprintf("%s/index", sha1FileDir))
	if err != nil {
		return nil, err
	}
	header, err := NewCacheHeaderFromBytes(bytes)
	if err != nil {
		return nil, err
	}
	return header.Entries, nil
}

type ActiveCache []*CacheEntry

func (ac ActiveCache) FindCacheEntryIndex(targetEntry *CacheEntry) int {
	for index, entry := range ac {
		if entry.Name == targetEntry.Name {
			return index
		}
	}
	return -1
}

func (ac ActiveCache) WriteCache(file *os.File) error {
	// SHA1ハッシュとる箇所の自信がない
	header := NewCacheHeader(1, ac)
	headerBytes := header.Bytes()
	file.Write(headerBytes)
	for _, entry := range ac {
		entryBytes := entry.Bytes()
		file.Write(entryBytes)
	}
	return nil
}
