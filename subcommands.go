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
	"strings"
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
	s := fmt.Sprintf("%s %s: unknown command\n", os.Args[0], subCmdName)
	s += fmt.Sprintf("Run '%s help' for usage.\n", os.Args[0])
	fmt.Fprint(c.topFlags.Output(), s)
	return nil
}

func (c *commands) printUsage() {
	var b strings.Builder
	b.WriteString("Usage:\n")
	if c.cmds != nil {
		fmt.Fprintf(&b, "\t%s <command> [arguments]\n\n", os.Args[0])
		b.WriteString("The commands are:\n")
		for _, cmd := range c.cmds {
			// Four spaces before the tab triggers good alignment
			// for both 4- and 8-space tab stops.
			fmt.Fprintf(&b, "\t%s    \t%s\n", cmd.Name(), cmd.Intro())
		}
	}
	fmt.Fprintf(&b, "\nUse \"%s help <command>\" for more information about a command.\n", os.Args[0])
	fmt.Fprint(c.topFlags.Output(), b.String())
}

var root = commands{
	topFlags: flag.CommandLine,
}

func Register(c Commander) {
	root.cmds = append(root.cmds, c)
}

func Execute() error {
	return root.Execute()
}
