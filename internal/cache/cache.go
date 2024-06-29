package cache

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"syscall"

	"github.com/marutaku/go-git/internal/cache/cachetime"
	"github.com/marutaku/go-git/internal/env"
	"github.com/marutaku/go-git/internal/hash"
	objectBuffer "github.com/marutaku/go-git/internal/objects"
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
	size := (62 + len(e.Name) + 8) & ^7
	bytes := make([]byte, size)
	binary.LittleEndian.PutUint32(bytes[0:], e.CTime.Sec)
	binary.LittleEndian.PutUint32(bytes[4:], e.CTime.NSec)
	binary.LittleEndian.PutUint32(bytes[8:], e.MTime.Sec)
	binary.LittleEndian.PutUint32(bytes[12:], e.MTime.NSec)
	binary.LittleEndian.PutUint32(bytes[16:], e.STDev)
	binary.LittleEndian.PutUint32(bytes[20:], e.STIno)
	binary.LittleEndian.PutUint32(bytes[24:], e.STMode)
	binary.LittleEndian.PutUint32(bytes[28:], e.STUid)
	binary.LittleEndian.PutUint32(bytes[32:], e.STGid)
	binary.LittleEndian.PutUint32(bytes[36:], e.STSize)
	copy(bytes[40:], e.Sha1)
	binary.LittleEndian.PutUint16(bytes[60:], e.NameLen)
	copy(bytes[62:], []byte(e.Name))
	return bytes
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
	sha1, err := hash.CalculateSha1HashFromFileStat(fileStat, fileContents)
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
	size := (62 + len(entry.Name) + 8) & ^7
	return entry, uint32(size)
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
