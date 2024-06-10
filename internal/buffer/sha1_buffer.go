package buffer

import (
	"fmt"
	"os"

	"github.com/marutaku/go-git/internal/env"
)

/*
* Blob, Tree, Commitの管理をする
 */

func getSha1FileName(sha1 []byte) string {
	sha1FileDirectory := env.GetSHA1FileDirectory()
	sha1Str := fmt.Sprintf("%x", sha1)
	return fmt.Sprintf("%s/%s/%s", sha1FileDirectory, sha1Str[:2], sha1Str[2:])
}

func WriteSha1Buffer(sha1 []byte, buffer []byte) error {
	fileName := getSha1FileName(sha1)
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	file.Write(buffer)
	return nil
}
