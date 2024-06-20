package utils

import (
	"bytes"
	"compress/zlib"
)

// Compress compresses the contents using zlib
func Compress(contents []byte) ([]byte, error) {
	var compressed bytes.Buffer
	writer := zlib.NewWriter(&compressed)
	_, err := writer.Write(contents)
	if err != nil {
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	return compressed.Bytes(), nil
}
