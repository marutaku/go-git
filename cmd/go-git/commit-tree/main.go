package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"os/user"

	"github.com/marutaku/go-git/internal/hash"
	"github.com/marutaku/go-git/internal/objects"
)

var MAX_PARENT = 16
var ORIG_OFFSET = 40

type CommitBuffer struct {
	buffer []byte
	offset int
}

func newBuffer() *CommitBuffer {
	return &CommitBuffer{
		buffer: make([]byte, ORIG_OFFSET),
		offset: ORIG_OFFSET,
	}
}

func (b *CommitBuffer) addBuffer(line string) {
	lineBytes := []byte(line)
	if len(lineBytes) >= len(b.buffer)-b.offset {
		appendBufferSize := (len(lineBytes) + 32767) &^ 32767
		b.buffer = append(b.buffer, make([]byte, appendBufferSize)...)
	}
	copy(b.buffer[b.offset:], lineBytes)
	b.offset += len(lineBytes)
}

func (b *CommitBuffer) finishBuffer(tag string) {
	start := objects.PrependInteger(b.buffer, b.offset-ORIG_OFFSET, ORIG_OFFSET)
	tagLen := len(tag)
	start -= tagLen
	copy(b.buffer[start:], []byte(tag))
	b.buffer = b.buffer[start:b.offset]
}

func getParentSha1s() ([][]byte, error) {
	// 以下のような形式で親コミットのSHA-1ハッシュ値が渡される
	// -p [parent sha1] -p [parent sha1] ...
	parentsCount := 0
	parentSha1s := make([][]byte, MAX_PARENT)
	for i := 2; i < len(os.Args); i += 2 {
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

func getRealCommitterName() (string, error) {
	var username string
	user, err := user.Current()
	if err != nil {
		return "", err
	}
	username = user.Username
	return username, nil
}

func getRealCommitterEmail() (string, error) {
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
	realCommitterName, err := getRealCommitterName()
	if err != nil {
		log.Fatal(err)
	}
	realCommitterEmail, err := getRealCommitterEmail()
	if err != nil {
		log.Fatal(err)
	}
	realCommitterDate, err := getCommitterDate()
	if err != nil {
		log.Fatal(err)
	}
	var committerName, committerEmail string
	if committerName = os.Getenv("GIT_COMMITTER_NAME"); committerName == "" {
		committerName = realCommitterName
	}
	if committerEmail = os.Getenv("GIT_COMMITTER_EMAIL"); committerEmail == "" {
		committerEmail = realCommitterEmail
	}
	committerDate := realCommitterDate
	if os.Getenv("GIT_COMMITTER_DATE") != "" {
		committerDate, err = time.Parse(time.RFC3339, os.Getenv("GIT_COMMITTER_DATE"))
		if err != nil {
			log.Fatal(err)
		}
	}
	// TODO: remove_special
	commitBuffer := newBuffer()
	commitBuffer.addBuffer(fmt.Sprintf("tree %s\n", treeSha1))
	for _, parentSha1 := range parentSha1s {
		commitBuffer.addBuffer(fmt.Sprintf("parent %s\n", parentSha1))
	}
	commitBuffer.addBuffer(fmt.Sprintf("author %s <%s> %d\n", realCommitterName, realCommitterEmail, realCommitterDate.Unix()))
	commitBuffer.addBuffer(fmt.Sprintf("committer %s <%s> %d\n", committerName, committerEmail, committerDate.Unix()))
	var comment string
	fmt.Scan(&comment)
	commitBuffer.addBuffer(comment)
	commitBuffer.finishBuffer("commit ")
	err = objects.WriteSha1File(commitBuffer.buffer)
	if err != nil {
		log.Fatal(err)
	}
}
