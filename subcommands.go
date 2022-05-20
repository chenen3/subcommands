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
	c.topFlags.Usage = c.printUsage
	// topFlags exit on parsing error
	c.topFlags.Parse(os.Args[1:])
	subCmdName := c.topFlags.Arg(0)
	if subCmdName == "" {
		c.topFlags.Usage()
		return nil
	}

	var help bool
	if subCmdName == "help" {
		help = true
		// help next subcommand
		subCmdName = c.topFlags.Arg(1)
		if subCmdName == "" {
			c.topFlags.Usage()
			return nil
		}
	}

	for _, cmd := range c.cmds {
		if cmd.Name() != subCmdName {
			continue
		}
		flags := flag.NewFlagSet(subCmdName, flag.ExitOnError)
		cmd.SetFlags(flags)
		if help {
			flags.Usage()
			return nil
		}
		flags.Parse(c.topFlags.Args()[1:])
		return cmd.Execute()
	}

	// unknown command
	s := fmt.Sprintf("unknown command: %s\n\n", subCmdName)
	s += fmt.Sprintf("See '%s help' for more information on the commands\n", os.Args[0])
	fmt.Fprint(c.topFlags.Output(), s)
	return nil
}

func (c *commands) printUsage() {
	s := fmt.Sprintf("Usage of %s:\n", os.Args[0])
	if c.cmds != nil {
		s += fmt.Sprintf("  %s <command>\n", os.Args[0])
		s += "\n"
		s += "Commands:\n"
		for _, cmd := range c.cmds {
			// Four spaces before the tab triggers good alignment
			// for both 4- and 8-space tab stops.
			s += fmt.Sprintf("  %s    \t%s\n", cmd.Name(), cmd.Intro())
		}
	}
	s += "\n"
	s += fmt.Sprintf("For more information on available commands and options see '%s help <command>'\n", os.Args[0])
	fmt.Fprint(c.topFlags.Output(), s)
}

var root = commands{
	topFlags: flag.CommandLine,
	cmds:     []Commander{new(helpCmd)},
}

func Register(c Commander) {
	root.cmds = append(root.cmds, c)
}

func Execute() error {
	return root.Execute()
}

type helpCmd struct{}

func (h *helpCmd) Name() string           { return "help" }
func (h *helpCmd) Intro() string          { return "print the help message" }
func (h *helpCmd) SetFlags(*flag.FlagSet) {}
func (h *helpCmd) Execute() error         { return nil }
