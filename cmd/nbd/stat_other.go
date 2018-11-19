// +build !linux,!darwin

package main

import (
	"os"

	"github.com/Merovius/nbd"
)

func blockSize(fi os.FileInfo) *nbd.BlockSizeConstraints {
	return nil
}
