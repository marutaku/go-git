package hash

import (
	"crypto/sha1"
	"encoding/hex"
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

func CalculateSha1HashFromFileFromByte(fileContent []byte) ([]byte, error) {
	h := sha1.New()
	h.Write(fileContent)
	sha1Bytes := h.Sum(nil)
	return sha1Bytes, nil
}

func GetSha1Hex(sha1Hash string) ([]byte, error) {
	bytes, err := hex.DecodeString(sha1Hash)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
