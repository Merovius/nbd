// +build linux darwin

package main

import (
	"os"

	"github.com/Merovius/nbd"
	"golang.org/x/sys/unix"
)

func blockSize(fi os.FileInfo) *nbd.BlockSizeConstraints {
	if st, ok := fi.Sys().(*unix.Stat_t); ok {
		if st.Blksize > 0xffff {
			return nil
		}
		return &nbd.BlockSizeConstraints{
			Min:       1,
			Preferred: uint32(st.Blksize),
			Max:       0xffff,
		}
	}
	return nil
}
