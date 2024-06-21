package hash

import (
	"crypto/sha1"
	"fmt"
	"io/fs"

	"github.com/marutaku/go-git/internal/utils"
)

func CalculateSha1HashFromFileStat(stat fs.FileInfo, fileContent []byte) ([]byte, error) {
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

func hexval(c int) int {
	if c >= 0x00 && c <= 0x09 {
		return c - '0'
	}
	if c >= 0x0a && c <= 0x0f {
		return c - 'a' + 10
	}
	if c >= 0x0A && c <= 0x0F {
		return c - 'A' + 10
	}
	return ^0
}

func CalculateSha1HashFromBytes(contents []byte) []byte {
	h := sha1.New()
	h.Write(contents)
	sha1Bytes := h.Sum(nil)
	return sha1Bytes
}

func GetSha1Hex(sha1Bytes []byte) bool {
	for i := 0; i < 20; i++ {

	}
}
