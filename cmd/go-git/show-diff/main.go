package main

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/marutaku/go-git/internal/buffer"
	"github.com/marutaku/go-git/internal/cache"
	"github.com/marutaku/go-git/internal/cache/cachetime"
)

var (
	MTIME_CHANGED = 0x0001
	CTIME_CHANGED = 0x0002
	OWNER_CHANGED = 0x0004
	MODE_CHANGED  = 0x0008
	INODE_CHANGED = 0x0010
	DATA_CHANGED  = 0x0020
)

func matchStat(entry *cache.CacheEntry, stat fs.FileInfo) int {
	changed := 0
	ctime := cachetime.NewCTimeFromStat(stat)
	mtime := cachetime.NewMTimeFromStat(stat)
	if ctime.Sec != entry.CTime.Sec || ctime.NSec != entry.CTime.NSec {
		changed |= CTIME_CHANGED
	}
	if mtime.Sec != entry.MTime.Sec || mtime.NSec != entry.MTime.NSec {
		changed |= MTIME_CHANGED
	}
	if stat.Sys().(*syscall.Stat_t).Uid != entry.STUid {
		changed |= OWNER_CHANGED
	}
	if stat.Mode() != fs.FileMode(entry.STMode) {
		changed |= MODE_CHANGED
	}
	if stat.Sys().(*syscall.Stat_t).Ino != uint64(entry.STIno) {
		changed |= INODE_CHANGED
	}
	if stat.Size() != int64(entry.STSize) {
		changed |= DATA_CHANGED
	}
	return changed
}

func showDifference(entry *cache.CacheEntry, oldContents []byte) error {
	executeCommand := fmt.Sprintf("diff -u - %s", entry.Name)
	cmd := exec.Command("/bin/bash", "-c", executeCommand)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	go func() {
		defer stdin.Close()
		io.WriteString(stdin, string(oldContents))
	}()
	// diffコマンドは差分があると終了ステータスが1になるため、エラーとして扱わない
	output, _ := cmd.CombinedOutput()
	fmt.Println(string(output))
	return nil
}

func main() {
	entries, err := cache.ReadCache()
	if err != nil {
		log.Fatal(err)
	}
	for _, entry := range entries {
		file, err := os.Open(entry.Name)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		fileStat, err := file.Stat()
		if err != nil {
			log.Fatal(err)
		}
		changed := matchStat(entry, fileStat)
		if changed == 0 {
			fmt.Printf("%s: ok\n", entry.Name)
			continue
		}
		fmt.Printf("%.*s: %02x", entry.NameLen, entry.Name, entry.Sha1)
		fmt.Print("\n")
		_, new, err := buffer.ReadSha1File(entry.Sha1)
		if err != nil {
			log.Fatal(err)
		}
		err = showDifference(entry, new)
		if err != nil {
			log.Fatal(err)
		}
	}
}
