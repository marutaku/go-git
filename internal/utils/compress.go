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
	byteReader := bytes.NewReader(compressed)
	zr, err := zlib.NewReader(byteReader)
	if err != nil {
		return nil, err
	}
	defer zr.Close()
	var decompressed []byte
	writer := bytes.NewBuffer(decompressed)
	if _, err := io.Copy(writer, zr); err != nil {
		return nil, err
	}
	return writer.Bytes(), nil
}
