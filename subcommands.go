/*
Copyright 2016 Google Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package subcommands implements a simple way for a single command to have many
// subcommands, each of which takes arguments and so forth.
package subcommands

import (
	"flag"
	"fmt"
	"os"
)

type Commander interface {
	// Name returns the name of the command.
	Name() string

	// Intro returns a short string (less than one line) describing the command.
	Intro() string

	// SetFlags adds the flags for this command to the specified set.
	SetFlags(*flag.FlagSet)

	// Execute executes the command.
	Execute() error
}

type commands struct {
	cmds     []Commander
	topFlags *flag.FlagSet
}

func (c *commands) Execute() error {
	// TODO: help
	// help:=c.topFlags.BoolVar("help", false, "print help message")
	err := c.topFlags.Parse(os.Args[1:])
	if err != nil {
		return err
	}

	subCmdName := c.topFlags.Arg(0)
	for _, cmd := range c.cmds {
		if cmd.Name() != subCmdName {
			continue
		}
		flags := flag.NewFlagSet(subCmdName, flag.ContinueOnError)
		help := flags.Bool("help", false, "print the help message")
		cmd.SetFlags(flags)
		if err := flags.Parse(c.topFlags.Args()[1:]); err != nil {
			return err
		}
		if *help {
			fmt.Fprintf(flags.Output(), "%s\n\n", cmd.Intro())
			flags.PrintDefaults()
			return nil
		}
		return cmd.Execute()
	}

	c.printUsage()
	return nil
}

func (c *commands) printUsage() {
	s := "Usage:\n"
	if c.topFlags.NFlag() > 0 {
		s += fmt.Sprintf("  %s [flags]\n", os.Args[0])
	}
	if c.cmds != nil {
		s += fmt.Sprintf("  %s [command]\n", os.Args[0])
		s += "\n"
		s += "Commands:\n"
		for _, cmd := range c.cmds {
			// Four spaces before the tab triggers good alignment
			// for both 4- and 8-space tab stops.
			s += fmt.Sprintf("  %s    \t%s\n", cmd.Name(), cmd.Intro())
		}
	}
	fmt.Fprint(c.topFlags.Output(), s)

	if c.topFlags.NFlag() > 0 {
		s += fmt.Sprintf("  %s [flags]\n", os.Args[0])
		c.topFlags.PrintDefaults()
	}
}

var root = commands{topFlags: flag.NewFlagSet(os.Args[0], flag.ContinueOnError)}

func Register(c Commander) {
	root.cmds = append(root.cmds, c)
}

func Execute() error {
	return root.Execute()
}
