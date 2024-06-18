package cache

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"errors"
)

const CACHE_SIGNATURE = "CRID" // 本当は"DIRC"だが、なぜか本家のindexを見ると"CRID"になっている...？

type CacheHeader struct {
	Signature string
	Version   uint32
	Entries   []*CacheEntry
}

func NewCacheHeader(version uint32, entries []*CacheEntry) *CacheHeader {
	return &CacheHeader{
		Signature: CACHE_SIGNATURE,
		Version:   1,
		Entries:   entries,
	}
}

func (h *CacheHeader) Verify(expectSha1 []byte) error {
	if h.Signature != CACHE_SIGNATURE {
		return errors.New("bad signature")
	}
	if h.Version != 1 {
		return errors.New("bad version")
	}
	if !bytes.Equal(h.Sha1Hash(), expectSha1) {
		return errors.New("bad header sha1")
	}
	return nil
}

func (h *CacheHeader) Sha1Hash() []byte {
	bytes := make([]byte, 0)
	bytes = append(bytes, h.Signature...)
	bytes = binary.LittleEndian.AppendUint32(bytes, h.Version)
	bytes = binary.LittleEndian.AppendUint32(bytes, uint32(len(h.Entries)))
	hash := sha1.New()
	hash.Write(bytes)
	for _, e := range h.Entries {
		hash.Write(e.Bytes())
	}
	return hash.Sum(nil)
}

func NewCacheHeaderFromBytes(bytes []byte) (*CacheHeader, error) {
	header := &CacheHeader{}
	header.Signature = string(bytes[:4])
	header.Version = binary.LittleEndian.Uint32(bytes[4:8])
	entryCount := binary.LittleEndian.Uint32(bytes[8:12])
	header.Entries = make([]*CacheEntry, entryCount)
	sha1HashFromByte := bytes[12:32]
	offset := uint32(32)
	for i := 0; i < int(entryCount); i++ {
		entry, size := NewCacheEntryFromBytes(bytes[offset:])
		header.Entries[i] = entry
		offset += size
	}
	if err := header.Verify(sha1HashFromByte); err != nil {
		return nil, err
	}
	return header, nil
}

func (h *CacheHeader) Bytes() []byte {
	bytes := make([]byte, 0)
	bytes = append(bytes, h.Signature...)
	bytes = binary.LittleEndian.AppendUint32(bytes, h.Version)
	bytes = binary.LittleEndian.AppendUint32(bytes, uint32(len(h.Entries)))
	bytes = append(bytes, h.Sha1Hash()...)
	return bytes
}
