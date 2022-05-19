package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/chenen3/subcommands"
)

type fooCmd struct {
	bar bool
}

func (f *fooCmd) Name() string  { return "foo" }
func (f *fooCmd) Intro() string { return "i am the subcommand foo" }

func (f *fooCmd) SetFlags(flags *flag.FlagSet) {
	flags.BoolVar(&f.bar, "bar", false, "")
}

func (f *fooCmd) Execute() error {
	fmt.Println("bar: ", f.bar)
	return nil
}

func main() {
	subcommands.Register(new(fooCmd))
	err := subcommands.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
