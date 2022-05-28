package subcommands

import (
	"flag"
	"os"
	"testing"
)

type fooCmd struct {
	setFlag  bool
	executed bool
}

func (f *fooCmd) Name() string  { return "foo" }
func (f *fooCmd) Intro() string { return "i am the subcommand foo" }

func (f *fooCmd) SetFlags(flags *flag.FlagSet) {
	flags.BoolVar(&f.setFlag, "bar", false, "")
}

func (f *fooCmd) Execute() error {
	f.executed = true
	return nil
}

func TestSubCmd(t *testing.T) {
	os.Args = []string{"subcmdtest", "foo", "-bar"}
	fc := new(fooCmd)
	Register(fc)
	err := Execute()
	if err != nil {
		t.Fatal(err)
	}
	if !fc.setFlag {
		t.Fatal("flag not setted")
	}
	if !fc.executed {
		t.Fatal("not executed")
	}
}
