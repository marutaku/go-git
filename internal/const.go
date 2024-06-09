package internal

import "os"

var DB_ENVIRONMENT_KEY = "SHA1_FILE_DIRECTORY"
var DEFAULT_DB_ENVIRONMENT = ".dircache/objects"

func GetSHA1FileDirectory() string {
	sha1FileDir := os.Getenv(DB_ENVIRONMENT_KEY)
	if sha1FileDir != "" {
		return sha1FileDir
	}
	return DEFAULT_DB_ENVIRONMENT
}
