package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/chenen3/subcommands"
)

type printCmd struct {
	b bool
}

func (f *printCmd) Name() string  { return "print" }
func (f *printCmd) Intro() string { return "Print args to stdout" }
func (f *printCmd) SetFlags(flags *flag.FlagSet) {
	flags.BoolVar(&f.b, "b", false, "bool output")
}

func (f *printCmd) Execute() error {
	fmt.Println("print", f.b)
	return nil
}

func main() {
	subcommands.Register(new(printCmd))
	err := subcommands.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
