package cache

import (
	"crypto/sha1"
	"fmt"

	"github.com/marutaku/go-git/internal/buffer"
	"github.com/marutaku/go-git/internal/utils"
)

func WriteSha1File(contents []byte) error {
	compressed, err := utils.Compress(contents)
	if err != nil {
		return err
	}
	h := sha1.New()
	h.Write(compressed)
	sha1Bytes := h.Sum(nil)
	fmt.Printf("%x\n", sha1Bytes)
	return buffer.WriteSha1Buffer(sha1Bytes, compressed)
}
