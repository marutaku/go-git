package cache

import (
	"bytes"
	"compress/zlib"
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

func (e *CacheEntry) IndexFd(fileContent string, stat fs.FileInfo) error {
	contents := []byte(fmt.Sprintf("blob %d", uint32(stat.Size())))
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
	return objectBuffer.WriteSha1Buffer(sha1Bytes, buffer.Bytes())
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

func NewCacheEntryFromBinary(indexFileBytes []byte) (*CacheEntry, uint32) {
	entry := &CacheEntry{}
	entry.CTime.Sec = binary.LittleEndian.Uint32(indexFileBytes)
	entry.CTime.NSec = binary.LittleEndian.Uint32(indexFileBytes[4:])
	entry.MTime.Sec = binary.LittleEndian.Uint32(indexFileBytes[8:])
	entry.MTime.NSec = binary.LittleEndian.Uint32(indexFileBytes[12:])
	entry.STDev = binary.LittleEndian.Uint32(indexFileBytes[16:])
	entry.STIno = binary.LittleEndian.Uint32(indexFileBytes[20:])
	entry.STMode = binary.LittleEndian.Uint32(indexFileBytes[24:])
	entry.STUid = binary.LittleEndian.Uint32(indexFileBytes[28:])
	entry.STGid = binary.LittleEndian.Uint32(indexFileBytes[32:])
	entry.STSize = binary.LittleEndian.Uint32(indexFileBytes[36:])
	entry.Sha1 = indexFileBytes[40:60]
	entry.NameLen = binary.LittleEndian.Uint16(indexFileBytes[60:])
	entry.Name = string(indexFileBytes[62 : 62+entry.NameLen])
	return entry, uint32(62 + entry.NameLen)
}

func ReadCache() (int, error) {
	sha1FileDir := env.GetSHA1FileDirectory()
	if _, err := os.Stat(sha1FileDir); os.IsExist(err) {
		return 0, errors.New("SHA1 file directory not found")
	}
	// TODO: Implement the rest of the function
	return 0, nil
}
