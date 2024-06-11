package cache

import "crypto/sha1"

const CACHE_SIGNATURE = 0x44495243

type CacheHeader struct {
	Signature uint
	Version   uint
	Entries   []*CacheEntry
}

func NewCacheHeader(version uint, entries []*CacheEntry) *CacheHeader {
	return &CacheHeader{
		Signature: CACHE_SIGNATURE,
		Version:   1,
		Entries:   entries,
	}
}

func (h *CacheHeader) Byte() []byte {
	bytes := make([]byte, 0)
	bytes = append(bytes, byte(h.Signature))
	bytes = append(bytes, byte(h.Version))
	hash := sha1.New()
	hash.Write(bytes)
	for _, e := range h.Entries {
		hash.Write(e.Sha1[:])
	}

}

func (h *CacheHeader) Export(entry *CacheEntry) []byte {
	bytes := make([]byte, 0)
	bytes = append(bytes, byte(h.Signature))
	bytes = append(bytes, byte(h.Version))
	hash := sha1.New()
	hash.Write(bytes)
	for _, e := range h.Entries {
		hash.Write(byte(e))
	}
}
