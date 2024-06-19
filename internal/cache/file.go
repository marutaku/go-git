package cache

import (
	"bytes"
	"compress/zlib"
)

func WriteSha1File(contents []byte) error {
	var buffer bytes.Buffer
	zWriter := zlib.NewWriter(&buffer)
	zWriter, err := zlib.NewWriterLevel(zWriter, zlib.BestCompression)
	if err != nil {
		return err
	}
	zWriter.Write(contents)
	sha1 := calculateSha1Hash()
}
