/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package Commands provides common definitions & functionality for a CLI tool
// subcommand implementations
package commands

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

// Handler a subcommand must provide
type Handler func(*Command, []string) int

// Command related data
type Command struct {
	name    string
	descr   string
	fset    *flag.FlagSet
	handler Handler
}

type cmdMap map[string]*Command

// Map of Commands keyed by a normalized command name.
// Provides a global table of all registered subcommands
type Map struct {
	cmdMap
	sync.RWMutex
}

func (c *Command) checkPtr() {
	if c == nil {
		log.Fatal("Nil command")
	}
}

func (c *Command) Name() string {
	c.checkPtr()
	return c.name
}

func (c *Command) Description() string {
	c.checkPtr()
	return c.descr
}

// Handle is Command's handler receiver that invokes the handles, each command
// implementation must provide
func (cmd *Command) Handle(arguments []string) int {
	cmd.checkPtr()
	if cmd.handler == nil {
		log.Fatal("Nil handler")
	}
	return cmd.handler(cmd, arguments)
}

// Flags returns the command's own FlagSet
func (cmd *Command) Flags() *flag.FlagSet {
	if cmd == nil {
		log.Fatal("Nil command")
	}
	return cmd.fset
}

// Usage prints the command's Usage help
func (cmd *Command) Usage() {
	if cmd != nil && cmd.fset != nil {
		cmd.fset.Usage()
	}
}

// Usage prints out all commands' usages shifted right by two tabs
func (cmds *Map) Usage() {
	if cmds == nil {
		return
	}
	cmds.RLock()
	defer cmds.RUnlock()

	// NOTE: std flag package uses os.Stderr for output
	// (https://golang.org/pkg/flag/#FlagSet.SetOutput), we preserve this
	// behavior
	out := os.Stderr
	for name, cmd := range cmds.cmdMap {
		// Underscore & bold command's name in the usage printout
		fmt.Fprintf(out, "\n\t\033[4m\033[1m%s\033[0m - %s\n", name, cmd.descr)
		f := cmd.Flags()
		b := bytes.NewBufferString("\t\t") // shift first string right by 2 tabs
		f.SetOutput(b)                     // replace std out stream with our buffer
		cmd.Usage()
		// Shift every new line right by two tabs (ASCII text only)
		bstr :=
			bytes.Replace(b.Bytes(), []byte{'\n'}, []byte{'\n', '\t', '\t'}, -1)
		os.Stderr.Write(bstr)
		f.SetOutput(out) // restore original output stream
		fmt.Fprintln(out)
	}
}

// Add adds a new subcommnd to the global registered commands table
func (cmds *Map) Add(name, descr string, handler Handler) *Command {
	if cmds == nil {
		log.Fatalf("Nil commands map")
		return nil // to allow switch to log.Print...
	}
	if handler == nil {
		log.Fatalf("Nil command handler")
		return nil
	}
	name = normalizeCmdStr(name)
	if len(name) == 0 {
		log.Fatalf("Empty command name")
		return nil
	}

	cmd := &Command{
		name,
		descr,
		flag.NewFlagSet(name, flag.ExitOnError),
		handler}

	cmds.Lock()
	defer cmds.Unlock()

	if cmds.cmdMap == nil {
		cmds.cmdMap = cmdMap{}
	}
	if _, ok := cmds.cmdMap[name]; ok {
		log.Fatalf("Command %s is already registered", name)
		return nil
	}
	cmds.cmdMap[name] = cmd

	return cmd
}

// Get finds & returns a command with the given name. Returns nil if not found
func (cmds *Map) Get(cmdName string) *Command {
	if cmds != nil {
		cmds.RLock()
		cmd, ok := cmds.cmdMap[normalizeCmdStr(cmdName)]
		cmds.RUnlock()
		if ok {
			return cmd
		}
	}
	return nil
}

// GetIdx finds & returns a command with cmdName and it's index in the command line. Returns (nil, idx) if not found
func (cmds *Map) GetIdx(cmdName string) (cmd *Command, idx int) {
	cmd = cmds.Get(cmdName)
	for i, arg := range os.Args {
		if arg == cmdName {
			idx = i
			return
		}
	}
	return
}

// GetCommand returns first command line command (non-flag argument) if any & its list of arguments
func (cmds *Map) GetCommand() (cmd *Command, cmdArgs []string) {
	var idx int
	cmd, idx = cmds.GetIdx(flag.Arg(0))
	if cmd == nil {
		return
	}
	return cmd, os.Args[idx+1:]
}

// HandleCommand parses command line & handles first found command with its arguments
// it returns command exit code with nil error if successful or error if command is not registered/invalid
// HandleCommand also prints help/usage into stdout if command is missing or is help/h
func (cmds *Map) HandleCommand() (exitCode int, err error) {
	cmd, args := cmds.GetCommand()
	if cmd == nil {
		cmdName := strings.ToLower(flag.Arg(0))
		if cmdName != "" && cmdName != "help" && cmdName != "h" {
			return 1, fmt.Errorf("invalid command: %s in %v", cmdName, flag.Args())
		}
		flag.Usage()
		return 1, nil
	}
	err = cmd.Flags().Parse(args)
	if err != nil {
		return 2, err
	}
	return cmd.Handle(args), nil
}

func normalizeCmdStr(cs string) string {
	return strings.ToLower(strings.TrimSpace(cs))
}
