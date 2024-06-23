package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"os/user"

	"github.com/marutaku/go-git/internal/hash"
)

var MAX_PARENT = 16

func getParentSha1s() ([][]byte, error) {
	// 以下のような形式で親コミットのSHA-1ハッシュ値が渡される
	// -p [parent sha1] -p [parent sha1] ...
	parentsCount := 0
	parentSha1s := make([][]byte, MAX_PARENT)
	for i := 2; i < len(parentSha1s); i += 2 {
		if os.Args[i] != "-p" {
			return nil, fmt.Errorf("invalid option: %s", os.Args[i])
		}
		sha1Bytes, err := hash.GetSha1Hex(os.Args[i+1])
		if err != nil {
			return nil, err
		}
		parentSha1s[parentsCount] = sha1Bytes
		parentsCount++
	}
	return parentSha1s[:parentsCount], nil
}

func getCommitterName() (string, error) {
	var username string
	if os.Getenv("COMMITTER_NAME") != "" {
		return os.Getenv("COMMITTER_NAME"), nil
	}
	user, err := user.Current()
	if err != nil {
		return "", err
	}
	username = user.Username
	return username, nil
}

func getCommitterEmail() (string, error) {
	if os.Getenv("COMMITTER_EMAIL") != "" {
		return os.Getenv("COMMITTER_EMAIL"), nil
	}
	user, err := user.Current()
	if err != nil {
		return "", err
	}
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s@%s", user.Username, hostname), err
}

func getCommitterDate() (time.Time, error) {
	if os.Getenv("COMMITTER_DATE") != "" {
		return time.Parse(time.RFC3339, os.Getenv("COMMITTER_DATE"))
	}
	return time.Now(), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("commit-tree <sha1> [-p <sha1>]* < changelog")
	}
	if len(os.Args) < 2 {
		log.Fatal("commit-tree <sha1> [-p <sha1>]* < changelog")
	}
	treeSha1, err := hash.GetSha1Hex(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	parentSha1s, err := getParentSha1s()
	if err != nil {
		log.Fatal(err)
	}
	if len(parentSha1s) == 0 {
		fmt.Printf("Committing initial tree %s\n", os.Args[1])
	}
	committerName, err := getCommitterName()
	if err != nil {
		log.Fatal(err)
	}
	committerEmail, err := getCommitterEmail()
	if err != nil {
		log.Fatal(err)
	}
	committerDate, err := getCommitterDate()
	if err != nil {
		log.Fatal(err)
	}

}
