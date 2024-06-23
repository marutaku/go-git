package cache

import (
	"fmt"

	"github.com/marutaku/go-git/internal/buffer"
	"github.com/marutaku/go-git/internal/hash"
	"github.com/marutaku/go-git/internal/utils"
)

func WriteSha1File(contents []byte) error {
	compressed, err := utils.Compress(contents)
	if err != nil {
		return err
	}
	sha1Bytes, err := hash.CalculateSha1HashFromFileFromByte(compressed)
	if err != nil {
		return err
	}
	fmt.Printf("%x\n", sha1Bytes)
	return buffer.WriteSha1Buffer(sha1Bytes, compressed)
}
