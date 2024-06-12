//go:build darwin
// +build darwin

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
		Sec:  uint32(fileStat.Sys().(*syscall.Stat_t).Ctimespec.Sec),
		NSec: uint32(fileStat.Sys().(*syscall.Stat_t).Ctimespec.Nsec),
	}
}

func NewMTimeFromStat(fileStat fs.FileInfo) *CacheTime {
	return &CacheTime{
		Sec:  uint32(fileStat.Sys().(*syscall.Stat_t).Mtimespec.Sec),
		NSec: uint32(fileStat.Sys().(*syscall.Stat_t).Mtimespec.Nsec),
	}
}
