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

func (p *printCmd) Name() string  { return "print" }
func (p *printCmd) Intro() string { return "print args to stdout" }
func (p *printCmd) SetFlags(flags *flag.FlagSet) {
	flags.BoolVar(&p.b, "b", false, "bool output")
}

func (p *printCmd) Execute() error {
	fmt.Println("print", p.b)
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
