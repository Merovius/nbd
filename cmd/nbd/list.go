// +build linux

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/Merovius/nbd/nbdnl"
	"github.com/google/subcommands"
)

func init() {
	commands = append(commands, &listCmd{})
}

type listCmd struct {
}

func (cmd *listCmd) Name() string {
	return "list"
}

func (cmd *listCmd) Synopsis() string {
	return "list NBD devices and their status"
}

func (cmd *listCmd) Usage() string {
	return `Usage: nbd list

List NBD devices and their status
`
}

func (cmd *listCmd) SetFlags(fs *flag.FlagSet) {
}

func (cmd *listCmd) Execute(ctx context.Context, fs *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	st, err := nbdnl.StatusAll()
	if err != nil {
		log.Println(err)
		return subcommands.ExitFailure
	}
	sort.Slice(st, func(i, j int) bool { return st[i].Index < st[j].Index })

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	fmt.Fprintf(w, "Device\tConnected\n")
	for _, s := range st {
		fmt.Fprintf(w, "/dev/nbd%d\t%v\n", s.Index, s.Connected)
	}
	w.Flush()
	return subcommands.ExitSuccess
}
