# subcommands #

[![GoDoc](https://godoc.org/github.com/chenen3/subcommands?status.svg)](https://godoc.org/github.com/chenen3/subcommands)  
Subcommands is a Go package that implements a simple way for a single command to
have many subcommands, each of which takes arguments and so forth.

## Usage ##

Set up a 'foo' subcommand:

```go
import "flag"

type fooCmd struct {
  bar bool
  cat bool
}

func (f *fooCmd) Name() string  { return "foo" }
func (f *fooCmd) Intro() string { return "i am the subcommand foo" }

func (f *fooCmd) SetFlags(flags *flag.FlagSet) {
  flags.BoolVar(&f.bar, "bar", false, "")
}

func (f *fooCmd) Execute() error {
  f.cat = true
  return nil
}
```

Register the subcommands and execute:

```go
func main() {
	subcommands.Register(new(fooCmd))
	err := subcommands.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
```
