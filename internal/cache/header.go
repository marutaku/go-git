package cache

import (
	"crypto/sha1"
	"encoding/binary"
)

const CACHE_SIGNATURE = "CRID" // 本当は"DIRC"だが、なぜか本家は"CRID"になっている...？

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

func NewCacheHeaderFromBytes(bytes []byte) *CacheHeader {
	header := &CacheHeader{}
	header.Signature = string(bytes[:4])
	header.Version = binary.LittleEndian.Uint32(bytes[4:8])
	entryCount := binary.LittleEndian.Uint32(bytes[8:12])
	header.Entries = make([]*CacheEntry, entryCount)
	// sha1ハッシュを読み飛ばす
	offset := uint32(12 + 20)
	for i := 0; i < int(entryCount); i++ {
		entry, size := NewCacheEntryFromBytes(bytes[offset:])
		header.Entries[i] = entry
		offset += size
	}
	return header
}

func (h *CacheHeader) Bytes() []byte {
	bytes := make([]byte, 0)
	bytes = append(bytes, h.Signature...)
	bytes = binary.LittleEndian.AppendUint32(bytes, h.Version)
	bytes = binary.LittleEndian.AppendUint32(bytes, uint32(len(h.Entries)))
	hash := sha1.New()
	hash.Write(bytes)
	for _, e := range h.Entries {
		hash.Write(e.Bytes())
	}
	totalHash := hash.Sum(nil)
	bytes = append(bytes, totalHash...)
	return bytes
}
