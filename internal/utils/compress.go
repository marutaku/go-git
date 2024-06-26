package utils

import (
	"bytes"
	"compress/zlib"
	"io"
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

func Decompress(compressed []byte) ([]byte, error) {
	reader := bytes.NewReader(compressed)
	zr, err := zlib.NewReader(reader)
	if err != nil {
		return nil, err
	}
	defer zr.Close()
	var decompressed bytes.Buffer
	if _, err := io.Copy(&decompressed, zr); err != nil {
		return nil, err
	}

	return decompressed.Bytes(), nil
}
