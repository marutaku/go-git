//go:build linux
// +build linux

// darwin.goと同じ内容
// https://github.com/golang/go/issues/29393

package cachetime

import (
	"io/fs"
	"syscall"
)

type CacheTime struct {
	Sec  uint32
	NSec uint32
}

func NewCTimeFromStat(fileStat fs.FileInfo) *CacheTime {
	return &CacheTime{
		Sec:  uint32(fileStat.Sys().(*syscall.Stat_t).Ctim.Sec),
		NSec: uint32(fileStat.Sys().(*syscall.Stat_t).Ctim.Nsec),
	}
}

func NewMTimeFromStat(fileStat fs.FileInfo) *CacheTime {
	return &CacheTime{
		Sec:  uint32(fileStat.Sys().(*syscall.Stat_t).Mtim.Sec),
		NSec: uint32(fileStat.Sys().(*syscall.Stat_t).Mtim.Nsec),
	}
}
