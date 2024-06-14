package cache

import (
	"crypto/sha1"
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

func (h *CacheHeader) Bytes() []byte {
	bytes := make([]byte, 0)
	bytes = append(bytes, h.Signature...)
	bytes = append(bytes, byte(h.Version))
	hash := sha1.New()
	hash.Write(bytes)
	for _, e := range h.Entries {
		hash.Write(e.Bytes())
	}
	totalHash := hash.Sum(nil)
	bytes = append(bytes, totalHash...)
	return bytes
}
