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

type loCmd struct {
	addr string
	unix bool
}

func (cmd *loCmd) Name() string {
	return "lo"
}

func (cmd *loCmd) Synopsis() string {
	return "lo NBD devices and their status"
}

func (cmd *loCmd) Usage() string {
	return `Usage: nbd lo <file>

Serve a file as a network block device.
`
}

func (cmd *loCmd) SetFlags(fs *flag.FlagSet) {
	fs.StringVar(&cmd.addr, "addr", "localhost:10809", "Address to listen on")
	fs.BoolVar(&cmd.unix, "unix", false, "Treat -addr as a unix domain socket")
}

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
		<-ch
		d.crash()
		signal.Stop(ch)
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

func (c *crashable) crash() {
	atomic.CompareAndSwapUint32(&c.crashed, 0, 1)
}

func (c *crashable) WriteAt(p []byte, offset int64) (n int, err error) {
	if atomic.LoadUint32(&c.crashed) != 0 {
		return 0, nbd.Errorf(nbd.EIO, "crash simulated")
	}
	return c.Device.WriteAt(p, offset)
}
