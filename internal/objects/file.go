package objects

import (
	"bytes"
	"fmt"
	"os"

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
	return WriteSha1Buffer(sha1Bytes, compressed)
}

func ReadSha1File(sha1 []byte) (string, []byte, error) {
	var nodeType string
	var bodySize int
	fileName := GetSha1FileName(sha1)
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0644)
	fileBuffer := make([]byte, 8192)
	if err != nil {
		return "", nil, err
	}
	defer file.Close()
	size, err := file.Read(fileBuffer)
	if err != nil {
		return "", nil, err
	}
	decompressed, err := utils.Decompress(fileBuffer[:size])
	if err != nil {
		return "", nil, err
	}
	splittedBytes := bytes.Split(decompressed, []byte{0})
	headerBytes := splittedBytes[0]
	header := string(headerBytes)
	fmt.Sscanf(header, "%s %d", &nodeType, &bodySize)
	return nodeType, bytes.Join(splittedBytes[1:], []byte{0}), nil
}
