// +build linux

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync/atomic"

	"github.com/Merovius/nbd"
	"github.com/google/subcommands"
	"golang.org/x/sys/unix"
)

func init() {
	commands = append(commands, &loCmd{})
}

type loCmd struct{}

func (cmd *loCmd) Name() string {
	return "lo"
}

func (cmd *loCmd) Synopsis() string {
	return "Provide file locally as a block device"
}

func (cmd *loCmd) Usage() string {
	return `Usage: nbd lo <file>

Provide file locally as a block device. An NBD device node will be chosen automatically and the path of that device printed to stdout.

As a special feature, you can toggle write-only mode by sending a SIGUSR1. In
write-only mode, all write-requests are denied with a EPERM. This is useful for
testing crash-resilience of an application on a given filesystem. You can
create a virtual block device with a filesystem of your choice and have the
application under test write to it. When you want to simulate a crash, you send
a SIGUSR1 and unmount the device. You then send another SIGUSR1 and remount the
filesystem to check whether invariants of the application survived the "crash".
`
}

func (cmd *loCmd) SetFlags(fs *flag.FlagSet) {}

func (cmd *loCmd) Execute(ctx context.Context, fs *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if fs.NArg() != 1 {
		log.Print(cmd.Usage())
		return subcommands.ExitUsageError
	}

	f, err := os.OpenFile(fs.Arg(0), os.O_RDWR, 0)
	if err != nil {
		log.Println(err)
		return subcommands.ExitFailure
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		log.Println(err)
		return subcommands.ExitFailure
	}
	log.Println(fi.Size())

	d := &crashable{Device: f}
	ch := make(chan os.Signal)
	signal.Notify(ch, unix.SIGUSR1)
	go func() {
		for range ch {
			d.toggleCrash()
		}
	}()

	idx, wait, err := nbd.Loopback(ctx, d, uint64(fi.Size()))
	if err != nil {
		log.Println(err)
		return subcommands.ExitFailure
	}
	fmt.Printf("Connected to /dev/nbd%d\n", idx)
	if err := wait(); err != nil {
		log.Println(err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}

type crashable struct {
	nbd.Device
	crashed uint32
}

func (c *crashable) toggleCrash() {
	if atomic.AddUint32(&c.crashed, 1<<31) == 0 {
		log.Println("SIGUSR1 received, device is read-write")
	} else {
		log.Println("SIGUSR1 received, device is read-only")
	}
}

func (c *crashable) WriteAt(p []byte, offset int64) (n int, err error) {
	if atomic.LoadUint32(&c.crashed) != 0 {
		return 0, nbd.Errorf(nbd.EPERM, "write-only")
	}
	return c.Device.WriteAt(p, offset)
}
