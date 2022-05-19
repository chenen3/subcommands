package subcommands

import (
	"flag"
	"os"
	"testing"
)

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

func TestSubCmd(t *testing.T) {
	os.Args = []string{"subcmdtest", "foo", "-bar"}
	fc := new(fooCmd)
	Register(fc)
	err := Execute()
	if err != nil {
		t.Fatal(err)
	}
	if !fc.bar {
		t.Fatalf("want bar: %v, got %v", true, fc.bar)
	}
	if !fc.cat {
		t.Fatalf("want cat: %v, got %v", true, fc.cat)
	}
}
