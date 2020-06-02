/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Command Line Tool to create & manage Operators, ACLs and Certificates
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/tools/accessc/handlers"
)

func main() {
	flag.Parse()
	plugin.LoadAllPluginsFatalOnError(&plugin.DefaultOrchestratorPluginLoader{})

	// Init help for all commands
	flag.Usage = func() {
		cmd := os.Args[0]
		fmt.Printf(
			"\nUsage: \033[1m%s [GENERAL OPTIONS] command [COMMAND OPTIONS]\033[0m\n\n",
			filepath.Base(cmd))
		flag.PrintDefaults()
		fmt.Println("\nCommands:")
		handlers.CommandRegistry.Usage()
	}

	cmd, args := handlers.CommandRegistry.GetCommand()
	if cmd == nil {
		cmdName := strings.ToLower(flag.Arg(0))
		if cmdName != "" && cmdName != "help" && cmdName != "h" {
			fmt.Println("\nInvalid Command: ", cmdName)
		}
		flag.Usage()
		os.Exit(1)
	}
	cmd.Flags().Parse(args)
	os.Exit(cmd.Handle(args))
}
