package buffer

import (
	"fmt"
	"os"

	"github.com/marutaku/go-git/internal/env"
)

/*
* Blob, Tree, Commitの管理をする
 */

func GetSha1FileName(sha1 []byte) string {
	sha1FileDirectory := env.GetSHA1FileDirectory()
	sha1Str := fmt.Sprintf("%x", sha1)
	return fmt.Sprintf("%s/%s/%s", sha1FileDirectory, sha1Str[:2], sha1Str[2:])
}

func WriteSha1Buffer(sha1 []byte, buffer []byte) error {
	fileName := GetSha1FileName(sha1)
	fmt.Println(fileName)
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		if os.IsExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()
	file.Write(buffer)
	return nil
}

func PrependInteger(buffer []byte, value int, offset int) int {
	// 本当は'\0'が入る
	buffer[offset-1] = '\x00'
	offset--
	for value > 0 {
		buffer[offset] = '0' + byte(value%10)
		value /= 10
		offset--
	}
	return offset
}
