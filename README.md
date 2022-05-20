# subcommands #

[![GoDoc](https://godoc.org/github.com/chenen3/subcommands?status.svg)](https://godoc.org/github.com/chenen3/subcommands)  
Subcommands is a Go package that implements a simple way for a single command to
have many subcommands, each of which takes arguments and so forth.

## Usage ##

Set up a 'foo' subcommand:

```go
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

```

Register the subcommands and execute:

```go
func main() {
	subcommands.Register(new(printCmd))
	err := subcommands.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
```
